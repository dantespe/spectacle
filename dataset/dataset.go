// Dataset contains logic for creating a collection of data.
package dataset

import (
	"fmt"

	"github.com/dantespe/spectacle/db"
)

// Dataset contains all logic for managing data in Spectacle.
type Dataset struct {
	// DatasetId of the dataset.
	DatasetId int64 `json:"datasetId"`

	// DisplayName of Dataset.
	DisplayName string `json:"displayName"`

	// NumRecords in the dataset.
	NumRecords int64 `json:"numRecords"`

	// HeadersSet
	HeadersSet bool `json:"headersSet"`

	MinRecordId int64 `json:"-"`

	MaxRecordId int64 `json:"-"`

	eng *db.Engine
}

// Option for creating new Datasets.
type Option func(*Dataset)

// New creates a new Dataset with the given options.
func New(eng *db.Engine, opts ...Option) (*Dataset, error) {
	if eng == nil {
		return nil, fmt.Errorf("eng must be non-nil")
	}

	// Create Dataset based on our options
	ds := &Dataset{
		HeadersSet: false,
		NumRecords: 0,
		eng:        eng,
	}
	for _, o := range opts {
		o(ds)
	}

	// Insert Dataset into DB and update the ds Id
	err := eng.DatabaseHandle.QueryRow("INSERT INTO Datasets (DisplayName, HeadersSet, NumRecords) VALUES ($1, 0, 0) RETURNING DatasetId", ds.DisplayName).Scan(&ds.DatasetId)
	if err != nil {
		return nil, fmt.Errorf("failed to create Dataset with error: %v", err)
	}

	// Update DisplayName to untitled-<datatsetId> if unset
	if ds.DisplayName == "" {
		ds.DisplayName = fmt.Sprintf("untitled-%d", ds.DatasetId)
		stmt, err := eng.DatabaseHandle.Prepare("UPDATE Datasets SET DisplayName = $1 WHERE DatasetId = $2")
		if err != nil {
			return nil, fmt.Errorf("failed to create Datasets prepared statement with error: %v", err)
		}
		_, err = stmt.Exec(ds.DisplayName, ds.DatasetId)
		if err != nil {
			return nil, fmt.Errorf("failed to insert into Datasets table with error: %v", err)
		}
	}
	return ds, nil
}

// Returns an Option with the DisplayName set.
func WithDisplayName(displayName string) Option {
	return func(ds *Dataset) {
		ds.DisplayName = displayName
	}
}

func GetDatasetFromId(eng *db.Engine, datasetId int64) (*Dataset, error) {
	if eng == nil {
		return nil, fmt.Errorf("eng must be non-nil")
	}

	// Get Dataset
	rows, err := eng.DatabaseHandle.Query("SELECT DisplayName, HeadersSet, NumRecords, MinRecordId, MaxRecordId FROM Datasets WHERE DatasetId = $1", datasetId)
	if err != nil {
		return nil, fmt.Errorf("failed to query for dataset with error: %v", err)
	}
	defer rows.Close()
	ds := &Dataset{
		DatasetId: datasetId,
		eng:       eng,
	}
	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(&ds.DisplayName, &ds.HeadersSet, &ds.NumRecords, &ds.MinRecordId, &ds.MaxRecordId); err != nil {
		// 404: We did not find the dataset given datasetId
		return nil, nil
	}

	return ds, nil
}

func TotalDatasets(eng *db.Engine) (int64, error) {
	if eng == nil {
		return 0, fmt.Errorf("eng must be non-nil")
	}

	var result int64
	row := eng.DatabaseHandle.QueryRow("SELECT COUNT(*) FROM Datasets")
	if err := row.Scan(&result); err != nil {
		return 0, fmt.Errorf("got error for COUNT(*) with error: %v", err)
	}
	return result, nil
}

func GetDatasets(eng *db.Engine, maxDatasets int64) ([]*Dataset, error) {
	if eng == nil {
		return nil, fmt.Errorf("eng must be non-nil")
	}
	if maxDatasets <= 0 {
		maxDatasets = 100
	}
	rows, err := eng.DatabaseHandle.Query("SELECT DatasetId, DisplayName, HeadersSet, NumRecords FROM Datasets ORDER BY DatasetId LIMIT $1", maxDatasets)
	if err != nil {
		return nil, fmt.Errorf("failed to query for datasetId with error: %v", err)
	}
	defer rows.Close()

	results := make([]*Dataset, 0)
	for rows.Next() {
		ds := &Dataset{
			eng: eng,
		}
		if err := rows.Scan(&ds.DatasetId, &ds.DisplayName, &ds.HeadersSet, &ds.NumRecords); err != nil {
			return nil, fmt.Errorf("failed to Scan(DatasetId, DisplayName, HeadersSet, NumRecords) for dataset with error: %v", err)
		}
		results = append(results, ds)
	}
	return results, nil
}

func (d *Dataset) SetHeaders(headers bool) error {
	if d.eng == nil {
		return fmt.Errorf("eng must be non-nil")
	}

	h := 1
	if !headers {
		h = 0
	}
	stmt, err := d.eng.DatabaseHandle.Prepare("UPDATE Datasets SET HeadersSet = $1 WHERE DatasetId = $2")
	if err != nil {
		return fmt.Errorf("failed to create update Datasets prepared statement with error: %v", err)
	}
	_, err = stmt.Exec(h, d.DatasetId)
	if err != nil {
		return fmt.Errorf("failed to update into Datasets table with error: %v", err)
	}
	d.HeadersSet = headers
	return nil
}

func (d *Dataset) UpdateNumRecords() error {
	// Update TotalNumRecords
	stmt, err := d.eng.DatabaseHandle.Prepare("UPDATE Datasets SET NumRecords = (SELECT COUNT(*) FROM RecordsProcessed WHERE DatasetId = $1) WHERE DatasetId = $1")
	if err != nil {
		return fmt.Errorf("failed to create dataset NumRecords prepared statement with error: %v", err)
	}
	_, err = stmt.Exec(d.DatasetId)
	if err != nil {
		return fmt.Errorf("failed to update dataset NumRecords with error: %v", err)
	}

	// Update Min Record
	stmt, err = d.eng.DatabaseHandle.Prepare("UPDATE Datasets SET MinRecordId = (SELECT MIN(RecordId) FROM RecordsProcessed WHERE DatasetId = $1) WHERE DatasetId = $1")
	if err != nil {
		return fmt.Errorf("failed to create dataset NumRecords prepared statement with error: %v", err)
	}
	_, err = stmt.Exec(d.DatasetId)
	if err != nil {
		return fmt.Errorf("failed to update dataset NumRecords with error: %v", err)
	}

	// Update Max Record
	stmt, err = d.eng.DatabaseHandle.Prepare("UPDATE Datasets SET MaxRecordId = (SELECT MAX(RecordId) FROM RecordsProcessed WHERE DatasetId = $1) WHERE DatasetId = $1")
	if err != nil {
		return fmt.Errorf("failed to create dataset NumRecords prepared statement with error: %v", err)
	}
	_, err = stmt.Exec(d.DatasetId)
	if err != nil {
		return fmt.Errorf("failed to update dataset NumRecords with error: %v", err)
	}

	// Update values
	if err := d.eng.DatabaseHandle.QueryRow("SELECT NumRecords, MinRecordId, MaxRecordId FROM Datasets WHERE DatasetId = $1", d.DatasetId).Scan(&d.NumRecords, &d.MinRecordId, &d.MaxRecordId); err != nil {
		return fmt.Errorf("failed to retrieve dataset NumRecords with error: %v", d.DatasetId)
	}
	return nil
}
