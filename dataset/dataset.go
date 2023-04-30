// Dataset contains logic for creating a collection of data.
package dataset

import (
	"sync"
)

// Dataset contains all logic for managing data in Spectacle.
type Dataset struct {
	// Id of the dataset.
	Id uint64

	// DisplayName of Datast.
	DisplayName string

	// Number of Records (Rows) in the dataset.
	NumRecords uint64

	// Maximum number of threads to run when importing data.
	maxThreads int

	// Lock
	mu sync.RWMutex
}

// Default returns an initialized Dataset.
func Default() (*Dataset, error) {
	return &Dataset{
		maxThreads: 100,
	}, nil
}

// NewWithId builds a Dataset with the provided Id.
func NewWithId(id uint64) (*Dataset, error) {
	ds, _ := Default()
	ds.Id = id
	return ds, nil
}

// Summary for the Dataset.
func (d *Dataset) Summary() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return map[string]interface{} {
		"datasetId": d.Id,
		"displayName": d.DisplayName,
		"numRecords": d.NumRecords,
	}
}