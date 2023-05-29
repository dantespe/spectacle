package manager

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/operation"
)

// Manager stores all useful things for Spectacle.
type Manager struct {
	ds  map[uint64]*dataset.Dataset
	ops map[uint64]*operation.Operation
	uds int
	mu  sync.RWMutex
}

// New creates a new Manager.
func New() (*Manager, error) {
	return &Manager{
		ds:  make(map[uint64]*dataset.Dataset),
		ops: make(map[uint64]*operation.Operation),
		// one-based indexing on DisplayName
		uds: 1,
	}, nil
}

// Status returns the status of the server.
func (m *Manager) Status() (int, *StatusResponse) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return http.StatusOK, &StatusResponse{
		NumRecords: len(m.ds),
		Status:     "HEALTHY",
		Code:       http.StatusOK,
	}
}

// CreateDataset atomically creates a dataset.
func (m *Manager) CreateDataset(req *CreateDatasetRequest) (int, *CreateDatasetResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO: Add Maxmium number of retries (10).
	// Attempt to get a Unique ID
	for {
		id := rand.Uint64()
		if _, ok := m.ds[id]; !ok {
			ds, err := dataset.New(
				dataset.WithId(id),
				dataset.WithDisplayName(req.DisplayName),
				dataset.WithHasHeaders(req.HasHeaders),
			)
			if err != nil {
				return http.StatusBadRequest, &CreateDatasetResponse{
					Message: err.Error(),
				}
			}
			// Set DisplayName
			m.uds = ds.SetUntitledDisplayName(m.uds)
			m.ds[id] = ds

			return http.StatusCreated, &CreateDatasetResponse{
				DatasetUrl:  fmt.Sprintf("/dataset/%d", id),
				DatasetId:   id,
				DisplayName: ds.DisplayName,
			}
		}
	}
}

func (m *Manager) GetDataset(req *GetDatasetRequest) (int, *GetDatasetResponse) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ds, ok := m.ds[req.DatasetId]
	if !ok {
		return http.StatusNotFound, &GetDatasetResponse{
			Message: fmt.Sprintf("failed to find dataset with id: %d", req.DatasetId),
			Code:    http.StatusNotFound,
		}
	}
	return http.StatusOK, &GetDatasetResponse{
		Dataset: ds.Copy(),
		Code:    http.StatusOK,
	}
}

func (m *Manager) ListDatasets(req *ListDatasetsRequest) (int, *ListDatasetsResponse) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	resp := &ListDatasetsResponse{
		Results:       []*dataset.Dataset{},
		TotalDatasets: 0,
	}
	for _, ds := range m.ds {
		if ds != nil {
			resp.Results = append(resp.Results, ds.Copy())
		}
	}
	resp.TotalDatasets = len(resp.Results)
	return http.StatusOK, resp
}

func (m *Manager) createUniqueOperation() *operation.Operation {
	for {
		id := rand.Uint64()
		if _, ok := m.ops[id]; !ok {
			op, err := operation.New(id)
			if err == nil {
				m.ops[id] = op
				return op
			}
		}
	}
}

func (m *Manager) processUploadDatasetRequest(req *UploadDatasetRequest, op *operation.Operation, ds *dataset.Dataset) {
	op.MarkRunning()
	err := ds.Upload(req.InputFile)
	if err != nil {
		op.MarkFailed(err.Error())
	} else {
		op.MarkCompleted()
	}
}

func (m *Manager) UploadDataset(req *UploadDatasetRequest) (int, *UploadDatasetResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Query DatasetId
	ds, ok := m.ds[req.DatasetId]
	if !ok {
		return http.StatusNotFound, &UploadDatasetResponse{
			Message: fmt.Sprintf("failed to find dataset: %d", req.DatasetId),
		}
	}

	// Create Operation
	op := m.createUniqueOperation()

	// Non-blocking uploadDatasetRequset
	go m.processUploadDatasetRequest(req, op, ds)

	// return Operation
	return 200, &UploadDatasetResponse{
		OperationUrl: fmt.Sprintf("/operation/%d", op.Id),
	}
}
