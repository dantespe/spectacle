package manager

import (
    "fmt"
    "math/rand"
    "sync"

    "github.com/dantespe/spectacle/dataset"
    "github.com/dantespe/spectacle/operation"
)

// Manager stores all useful things for Spectacle.
type Manager struct {
    ds map[uint64]*dataset.Dataset
    ops map[uint64]*operation.Operation
    mu sync.RWMutex
}

// New creates a new Manager.
func New() (*Manager, error) {
    return &Manager{
        ds: make(map[uint64]*dataset.Dataset),
        ops: make(map[uint64]*operation.Operation),
    }, nil
}

// Status returns the status of the server.
func (m *Manager) Status() *StatusResponse {
    m.mu.RLock()
    defer m.mu.RUnlock()

    return &StatusResponse{
        NumRecords: len(m.ds),
    }
}

// CreateDataset atomically creates a dataset.
func (m *Manager) CreateDataset(req *CreateDatasetRequest) *CreateDatasetResponse {
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
            )
            m.ds[id] = ds
            return &CreateDatasetResponse{
                DatasetUrl: fmt.Sprintf("/dataset/%d", id),
                DatasetId: id,
                DisplayName: ds.DisplayName,
                Error: err,
            }
        }
    }
}

func (m *Manager) GetDataset(req *GetDatasetRequest) *GetDatasetResponse {
    m.mu.RLock()
    defer m.mu.RUnlock()

    ds, ok := m.ds[req.DatasetId]
    if !ok {
        return &GetDatasetResponse{
            Error: fmt.Errorf("failed to find dataset"),
        }
    }
    return &GetDatasetResponse{
        Dataset: ds,
    }
}

func (m *Manager) ListDatasets(req *ListDatasetsRequest) *ListDatasetsResponse {
    m.mu.RLock()
    defer m.mu.RUnlock()

    resp := newListDatasetsResponse()
    for _, ds := range m.ds {
        if ds != nil {
            resp.Results = append(resp.Results, ds.Summary())
        }
    }
    return resp
}