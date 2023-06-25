package manager

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/dantespe/spectacle/cell"
	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/db"
	"github.com/dantespe/spectacle/header"
	"github.com/dantespe/spectacle/operation"
	"github.com/dantespe/spectacle/record"
)

const (
	OPERATION_STATE_UNKNOWN = "UNKNOWN"
	OPERATION_STATE_RUNNING = "RUNNING"
	OPERATION_STATE_SUCCESS = "SUCCESS"
	OPERATION_STATE_FAIL    = "FAIL"
)

// Manager stores all useful things for Spectacle.
type Manager struct {
	mu  sync.RWMutex
	eng *db.Engine
}

// New creates a new Manager.
func New() (*Manager, error) {
	eng, err := db.New(
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

func (m *Manager) processUpload(req *UploadDatasetRequest, op *operation.Operation, ds *dataset.Dataset) {
	if err := op.MarkRunning(); err != nil {
		log.Printf("MarkRunning failed with error: %v", err)
		return
	}

	headers, err := header.GetHeaders(m.eng, req.DatasetId)
	if err != nil {
		log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
		op.MarkFailed(fmt.Sprintf("Failed to load headers into memory: %v", err))
	}

	// 2) Process Each Row in the CSV
	reader := csv.NewReader(req.InputFile)
	reader.ReuseRecord = true
	for {
		rawRecord, err := reader.Read()

		// Unexpected Error
		if err != nil && err != io.EOF {
			log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
			op.MarkFailed(fmt.Sprintf("Got unexpected error during processing: %v", err))
		}

		// EOF, wait until all threads have finished
		if err == io.EOF {
			break
		}

		// Create Headers
		if !ds.HeadersSet {
			for _ = range rawRecord {
				header, err := header.New(m.eng, req.DatasetId)
				if err != nil {
					log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
					op.MarkFailed(fmt.Sprintf("Got unexpected error during processing: %v", err))
				}
				headers = append(headers, header)
			}

			if ds.SetHeaders(true); err != nil {
				log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
				op.MarkFailed(fmt.Sprintf("Got unexpected error setting headers: %v", err))
			}
		}

		record, err := record.New(m.eng, req.DatasetId)
		if err != nil {
			log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
			op.MarkFailed(fmt.Sprintf("Failed to create new record with error: %v", err))
		}

		idx := 0
		for _, rv := range rawRecord {
			// Extend Headers if needed
			if idx >= len(headers) {
				header, err := header.New(m.eng, req.DatasetId)
				if err != nil {
					log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
					op.MarkFailed(fmt.Sprintf("Got unexpected error during processing: %v", err))
				}
				headers = append(headers, header)
			}
			if _, err := cell.New(m.eng, record.RecordId, headers[idx].HeaderId, op.OperationId, rv); err != nil {
				log.Printf("/operation/%d failed, check the logs to see a detailed error", op.OperationId)
				op.MarkFailed(fmt.Sprintf("Failed to create a cell: %v", err))
			}
			idx++
		}
	}

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
