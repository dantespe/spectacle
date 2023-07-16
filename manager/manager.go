package manager

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/db"
	"github.com/dantespe/spectacle/header"
	"github.com/dantespe/spectacle/operation"
)

// Manager stores all useful things for Spectacle.
type Manager struct {
	mu  sync.RWMutex
	eng *db.Engine
}

// New creates a new Manager.
func New() (*Manager, error) {
	eng, err := db.New(
		db.WithDatabaseProvider(db.DatabaseProvider_POSTGRES),
		db.WithEnvironment(db.Environment_DEVELOPMENT),
	)
	if err != nil {
		return nil, err
	}
	return NewWithEngine(eng)
}

func NewWithEngine(eng *db.Engine) (*Manager, error) {
	if eng == nil {
		return nil, fmt.Errorf("eng must be non-nil")
	}
	return &Manager{
		eng: eng,
	}, nil
}

// Status returns the status of the server.
func (m *Manager) Status() (int, *StatusResponse) {
	return http.StatusOK, &StatusResponse{
		Status: "HEALTHY",
		Code:   http.StatusOK,
	}
}

// CreateDataset atomically creates a dataset.
func (m *Manager) CreateDataset(req *CreateDatasetRequest) (int, *CreateDatasetResponse) {
	ds, err := dataset.New(m.eng, dataset.WithDisplayName(req.DisplayName))
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, &CreateDatasetResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	return http.StatusCreated, &CreateDatasetResponse{
		DatasetUrl:  fmt.Sprintf("/dataset/%d", ds.DatasetId),
		DatasetId:   ds.DatasetId,
		DisplayName: ds.DisplayName,
		Code:        http.StatusCreated,
	}
}

func (m *Manager) GetDataset(req *GetDatasetRequest) (int, *GetDatasetResponse) {
	ds, err := dataset.GetDatasetFromId(m.eng, req.DatasetId)
	if err != nil {
		log.Printf("Query for Dataset failed with error: %v", err)
		return http.StatusInternalServerError, &GetDatasetResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	// We should 404, because dataset was not found and err was nil.
	if ds == nil {
		return http.StatusNotFound, &GetDatasetResponse{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("failed to find dataset with id: %d", req.DatasetId),
		}
	}

	return http.StatusOK, &GetDatasetResponse{
		Code:    http.StatusOK,
		Dataset: ds,
	}
}

func (m *Manager) ListDatasets(req *ListDatasetsRequest) (int, *ListDatasetsResponse) {
	resp := &ListDatasetsResponse{
		Results:       []*dataset.Dataset{},
		TotalDatasets: 0,
		Code:          http.StatusOK,
	}

	td, err := dataset.TotalDatasets(m.eng)
	if err != nil {
		log.Printf("Failed to get total number of datasets with error: %v", err)
		return http.StatusInternalServerError, &ListDatasetsResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}
	resp.TotalDatasets = td

	// Add Datasets to Result
	results, err := dataset.GetDatasets(m.eng, req.MaxDatasets)
	if err != nil {
		log.Printf("Failed to get datasets with error: %v", err)
		return http.StatusInternalServerError, &ListDatasetsResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}
	resp.Results = results

	return http.StatusOK, resp
}

func (m *Manager) uploadRecordCount(ds *dataset.Dataset, op *operation.Operation) {
	// Every 10 seconds, we update the number of records in the Dataset.
	// We will timeout after 24 hours. This is way more time than needed. We can process
	// about 500k cells/min.
	ticker := time.NewTicker(10 * time.Second)
	timeout := time.NewTicker(24 * time.Hour)
	cancel := make(chan bool)

	go func() {
		for {
			select {
			case <-timeout.C:
				return
			case <-ticker.C:
				ds.UpdateNumRecords()
			case <-cancel:
				return
			}
		}
	}()

	for !op.Complete() {
		time.Sleep(time.Minute)
	}
	cancel <- true
}

func (m *Manager) createOrGetHeaders(rd io.Reader, op *operation.Operation, ds *dataset.Dataset) ([]*header.Header, error) {
	// Load Headers from Database
	headers, err := header.GetHeaders(m.eng, ds.DatasetId)
	if err != nil {
		return nil, err
	}

	// Return if we got a least one back
	if len(headers) > 0 && ds.HeadersSet {
		return headers, nil
	}
	if len(headers) > 0 && !ds.HeadersSet {
		err := ds.SetHeaders(true)
		return headers, err
	}

	// Read Headers from the file since we haven't seen any yet
	reader := csv.NewReader(rd)
	rawRecord, err := reader.Read()

	// Unexpected Error
	if err != nil && err != io.EOF {
		return nil, err
	}

	// EOF, so we return empty header list
	if err == io.EOF {
		// no headers
		return headers, nil
	}

	// Create Headers from input file
	tx, err := db.NewTx(m.eng, "headers", "datasetid", "displayname", "valuetype")
	if err != nil {
		return nil, err
	}
	for _, dn := range rawRecord {
		if err := tx.Exec(ds.DatasetId, dn, header.ValueType_RAW); err != nil {
			return nil, err
		}
	}
	if err := tx.Close(); err != nil {
		return nil, err
	}

	if err := ds.SetHeaders(true); err != nil {
		return nil, err
	}

	headers, err = header.GetHeaders(m.eng, ds.DatasetId)
	if err != nil {
		return nil, err
	}

	// Set Column Index Order
	for i, h := range headers {
		if err := h.SetColumnIndex(int64(i)); err != nil {
			return nil, err
		}
	}

	return headers, nil
}

func (m *Manager) createRecords(rd io.Reader, op *operation.Operation, ds *dataset.Dataset) error {
	// Create Records Transaction
	tx, err := db.NewTx(m.eng, "records", "operationid", "datasetid")
	if err != nil {
		return fmt.Errorf("failed to create record tx with err: %v", err)
	}

	// Create Each Record
	reader := csv.NewReader(rd)
	reader.LazyQuotes = true
	reader.ReuseRecord = true
	for {
		_, err := reader.Read()
		// Unexpected Error
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read record with err: %v", err)
		}

		// EOF, wait until all threads have finished
		if err == io.EOF {
			break
		}

		// Create Record
		if err := tx.Exec(op.OperationId, ds.DatasetId); err != nil {
			return fmt.Errorf("faield to Exec(op, ds) with err: %v", err)
		}
	}
	return tx.Close()
}

func (m *Manager) createCells(rd io.Reader, op *operation.Operation, ds *dataset.Dataset, headers []*header.Header) error {
	// Create RecordsProcessed Tx
	rtx, err := db.NewTx(m.eng, "recordsprocessed", "recordid", "datasetid")
	if err != nil {
		return err
	}

	// Create Cells Tx
	ctx, err := db.NewTx(m.eng, "cells", "recordid", "headerid", "operationid", "rawvalue")
	if err != nil {
		return err
	}

	// Get RecordIds
	records, err := m.eng.DatabaseHandle.Query("SELECT RecordId FROM Records WHERE OperationId = $1 ORDER BY RecordId", op.OperationId)
	if err != nil {
		return err
	}
	defer records.Close()

	reader := csv.NewReader(rd)
	reader.LazyQuotes = true
	reader.ReuseRecord = true
	first := true

	// For each record, we go through the CSV and create a cell for each
	// column in the row. Associate the correct foreign keys and then mark
	// the record (row) as processed.
	for records.Next() {
		var recordId int64
		if err := records.Scan(&recordId); err != nil {
			return err
		}

		// Read the row
		rawRecord, err := reader.Read()
		// Unexpected Error
		if err != nil && err != io.EOF {
			return err
		}
		// EOF
		if err == io.EOF {
			break
		}
		// First should skip
		if first {
			first = false
			continue
		}

		headerIdx := 0
		for _, rv := range rawRecord {
			// Extend Headers if needed
			if headerIdx >= len(headers) {
				header, err := header.New(m.eng, ds.DatasetId)
				if err != nil {
					return err
				}
				headers = append(headers, header)
			}
			// Create Cell for (row, col)
			if err := ctx.Exec(recordId, headers[headerIdx].HeaderId, op.OperationId, rv); err != nil {
				return err
			}
			headerIdx++
		}
		// Mark Record as Processed
		if err := rtx.Exec(recordId, ds.DatasetId); err != nil {
			return err
		}
	}

	if err := ctx.Close(); err != nil {
		return err
	}
	if err := rtx.Close(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) processUpload(req *UploadDatasetRequest, op *operation.Operation, ds *dataset.Dataset) {
	// Mark Operation as Running
	if err := op.MarkRunning(); err != nil {
		log.Printf("MarkRunning failed with error: %v", err)
		return
	}

	// Starts a background process to update the number of records in the dataset.
	go m.uploadRecordCount(ds, op)

	// Store Request File to Disk
	tmp, err := os.CreateTemp("", "spec_import")
	if err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to create temp file with error: %v", err))
		return
	}
	if _, err = io.Copy(tmp, req.InputFile); err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to copy to temp file with error: %v", err))
		return
	}
	defer os.Remove(tmp.Name())

	// Create Headers
	log.Printf("Creating Headers for operation: %d", op.OperationId)
	tf, err := os.Open(tmp.Name())
	if err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to copy to temp file for headers with error: %v", err))
		return
	}
	headers, err := m.createOrGetHeaders(tf, op, ds)
	if err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to load headers into memory: %v", err))
		return
	}

	// Create Records
	log.Printf("Creating Records for operation: %d", op.OperationId)
	rf, err := os.Open(tmp.Name())
	if err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to copy to temp file for records with error: %v", err))
		return
	}
	if err := m.createRecords(rf, op, ds); err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to create records: %v", err))
		return
	}

	// Create Cells
	log.Printf("Creating Cells for operation: %d", op.OperationId)
	cf, err := os.Open(tmp.Name())
	if err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to copy to temp file for cells with error: %v", err))
		return
	}
	if err := m.createCells(cf, op, ds, headers); err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to create cells: %v", err))
		return
	}
	log.Printf("Finishing operation: %d", op.OperationId)

	ds.UpdateNumRecords()
	op.MarkSuccess()
}

func (m *Manager) UploadDataset(req *UploadDatasetRequest) (int, *UploadDatasetResponse) {
	ds, err := dataset.GetDatasetFromId(m.eng, req.DatasetId)
	if err != nil {
		log.Printf("Query for Dataset failed with error: %v", err)
		return http.StatusInternalServerError, &UploadDatasetResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	if ds == nil {
		return http.StatusNotFound, &UploadDatasetResponse{
			Message: fmt.Sprintf("failed to find dataset with id: %d", req.DatasetId),
			Code:    http.StatusNotFound,
		}
	}

	op, err := operation.New(m.eng)
	if err != nil {
		log.Printf("Failed to build create operation statement with error: %v", err)
		return http.StatusInternalServerError, &UploadDatasetResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	go m.processUpload(req, op, ds)

	// Return Operation
	return http.StatusOK, &UploadDatasetResponse{
		OperationUrl: fmt.Sprintf("/operation/%d", op.OperationId),
		Code:         http.StatusOK,
	}
}

func (m *Manager) GetHeaders(req *GetHeadersRequest) (int, *GetHeadersResponse) {
	ds, err := dataset.GetDatasetFromId(m.eng, req.DatasetId)
	if err != nil {
		log.Printf("Query for Dataset failed with error: %v", err)
		return http.StatusInternalServerError, &GetHeadersResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	// We should 404, because dataset was not found and err was nil.
	if ds == nil {
		return http.StatusNotFound, &GetHeadersResponse{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("failed to find dataset with id: %d", req.DatasetId),
		}
	}

	headers, err := header.GetHeaders(m.eng, req.DatasetId)
	if err != nil {
		log.Printf("Failed to get headers with err: %v", err)
		return http.StatusInternalServerError, &GetHeadersResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	return http.StatusOK, &GetHeadersResponse{
		Headers: headers,
		Code:    http.StatusOK,
	}
}

func convertHeadersToString(headers []*header.Header) (string, error) {
	if headers == nil || len(headers) == 0 {
		return "", fmt.Errorf("can't convert headers for nil or empty list")
	}

	if len(headers) == 1 {
		return fmt.Sprintf("HeaderId = %d", headers[0].HeaderId), nil
	}

	result := "HeaderId IN ("
	for i, h := range headers {
		result += fmt.Sprintf("%d", h.HeaderId)
		// Append , if this is not the last header
		if i+1 < len(headers) {
			result += ","
		}
	}
	result += ")"
	return result, nil
}

func (m *Manager) GetData(req *DataRequest) (int, *DataResponse) {
	// Query for Dataset
	ds, err := dataset.GetDatasetFromId(m.eng, req.DatasetId)
	if err != nil {
		log.Printf("Query for Dataset failed with error: %v", err)
		return http.StatusInternalServerError, &DataResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	// We should 404, because dataset was not found and err was nil.
	if ds == nil {
		return http.StatusNotFound, &DataResponse{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("failed to find dataset with id: %d", req.DatasetId),
		}
	}

	// Get Headers
	headers, err := header.GetHeaders(m.eng, req.DatasetId)
	if err != nil {
		log.Printf("failed to get headers with err: %v", err)
		return http.StatusInternalServerError, &DataResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	// Exclude headers not included in req.Headers
	hasExclusions := false
	if len(req.Headers) > 0 {
		s := make(map[int64]bool, 0)
		for _, id := range req.Headers {
			s[id] = true
		}
		tmp := make([]*header.Header, 0)
		for _, h := range headers {
			if _, exists := s[h.HeaderId]; exists {
				hasExclusions = true
				tmp = append(tmp, h)
			}
		}
		headers = tmp
	}

	resp := &DataResponse{
		Code:    http.StatusOK,
		Results: make([]*ResultSet, 0),
	}

	if len(headers) == 0 {
		return http.StatusOK, resp
	}
	hs, err := convertHeadersToString(headers)
	if err != nil {
		log.Printf("failed to get headers string with err: %v", err)
		return http.StatusInternalServerError, &DataResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}
	hm := make(map[int64]*header.Header)
	for _, h := range headers {
		hm[h.HeaderId] = h
	}

	// Get MaxRecordId for this Block
	row := m.eng.DatabaseHandle.QueryRow("SELECT MAX(tmp.RecordId) FROM (SELECT RecordId FROM Records WHERE RecordId > $1 AND DatasetId = $2 LIMIT $3) AS tmp", req.LastRecordId, ds.DatasetId, req.MaxResults)
	var maxRecordId int64
	if err := row.Scan(&maxRecordId); err != nil {
		log.Printf("failed to get max record id with err: %v", err)
		return http.StatusInternalServerError, &DataResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	// Return Block of data
	limit := int64(len(headers)) * req.MaxResults
	q := fmt.Sprintf("SELECT RecordId, HeaderId, RawValue FROM Cells WHERE %s AND RecordId > %d AND RecordId <= %d ORDER BY HeaderId, RecordId  LIMIT %d", hs, req.LastRecordId, maxRecordId, limit)
	rows, err := m.eng.DatabaseHandle.Query(q)
	if err != nil {
		log.Printf("failed to build query for Cells with err: %v", err)
		return http.StatusInternalServerError, &DataResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	prevHeaderId := int64(-1)
	for rows.Next() {
		var recordId int64
		var headerId int64
		var rv string
		if err := rows.Scan(&recordId, &headerId, &rv); err != nil {
			log.Printf("failed to Scan(RecordId, HeaderId, RawValue) from Cells with err: %v", err)
			return http.StatusInternalServerError, &DataResponse{
				Message: "INTERNAL SERVER ERROR",
				Code:    http.StatusInternalServerError,
			}
		}

		header, ok := hm[headerId]
		if !ok {
			log.Printf("failed to get header displayName: %v", err)
			return http.StatusInternalServerError, &DataResponse{
				Message: "INTERNAL SERVER ERROR",
				Code:    http.StatusInternalServerError,
			}
		}
		if prevHeaderId < 0 || prevHeaderId != headerId {
			resp.Results = append(resp.Results, &ResultSet{
				HeaderId:    headerId,
				DisplayName: header.DisplayName,
				Data:        make([]string, 0),
			})
			prevHeaderId = headerId
		}
		resp.Results[len(resp.Results)-1].Data = append(resp.Results[len(resp.Results)-1].Data, rv)
	}

	var highestRecordId int64
	row = m.eng.DatabaseHandle.QueryRow("SELECT MAX(RecordId) FROM Records WHERE DatasetId = $1", ds.DatasetId)
	if err := row.Scan(&highestRecordId); err != nil {
		log.Printf("failed to get max recordid with error: %v", err)
		return http.StatusInternalServerError, &DataResponse{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
		}
	}

	// Populate the Next Page
	if highestRecordId > maxRecordId {
		baseUrl := fmt.Sprintf("/data/%d?recordid=%d", ds.DatasetId, maxRecordId+1)
		if hasExclusions {
			var headerIds []string
			for _, h := range headers {
				headerIds = append(headerIds, fmt.Sprintf("%d", h.HeaderId))
			}
			baseUrl += fmt.Sprintf("%s", strings.Join(headerIds, ","))
		}
		resp.Next = baseUrl
	}
	return http.StatusOK, resp
}
