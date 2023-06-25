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
		eng:        eng,
	}
	for _, o := range opts {
		o(ds)
	}

	// Insert Dataset into DB and update the ds Id
	stmt, err := eng.DatabaseHandle.Prepare("INSERT INTO Datasets(DisplayName, HeadersSet) VALUES(?, FALSE)")
	if err != nil {
		return nil, fmt.Errorf("failed to create Datasets prepared statement with error: %v", err)
	}
	res, err := stmt.Exec(ds.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into Datasets table with error: %v", err)
	}

	ds.DatasetId, err = res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to query DatasetId with error: %v", err)
	}

	// Update DisplayName to untitled-<datatsetId> if unset
	if ds.DisplayName == "" {
		ds.DisplayName = fmt.Sprintf("untitled-%d", ds.DatasetId)
		stmt, err := eng.DatabaseHandle.Prepare("UPDATE Datasets SET DisplayName = ? WHERE DatasetId = ?")
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
	rows, err := eng.DatabaseHandle.Query("SELECT DisplayName, HeadersSet FROM Datasets WHERE DatasetId = ?", datasetId)
	if err != nil {
		return nil, fmt.Errorf("failed to query for datasetId with error: %v", err)
	}
	defer rows.Close()

	// 404: We did not find the dataset with datasetId
	if !rows.Next() {
		return nil, nil
	}

	// Set DisplayName, HeadersSet
	ds := &Dataset{
		DatasetId: datasetId,
		eng:       eng,
	}
	if err := rows.Scan(&ds.DisplayName, &ds.HeadersSet); err != nil {
		return nil, fmt.Errorf("failed to Scan(displayName) for dataset with error: %v", err)
	}

	// Count Number of Records
	rows2, err := eng.DatabaseHandle.Query("SELECT COUNT(*) FROM Records WHERE DatasetId = ?", datasetId)
	defer rows2.Close()
	if rows2.Next() {
		if err := rows2.Scan(&ds.NumRecords); err != nil {
			return nil, fmt.Errorf("failed to Scan(NumRecords) for dataset with error: %v", err)
		}
	}
	return ds, nil
}

func TotalDatasets(eng *db.Engine) (int64, error) {
	if eng == nil {
		return 0, fmt.Errorf("eng must be non-nil")
	}

	row, err := eng.DatabaseHandle.Query("SELECT COUNT(*) FROM Datasets")
	defer row.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to get COUNT(*) from datasets with error: %v", err)
	}
	if !row.Next() {
		return 0, nil
	}

	var result int64
	err = row.Scan(&result)
	return result, err
}

func GetDatasets(eng *db.Engine, maxDatasets int64) ([]*Dataset, error) {
	if eng == nil {
		return nil, fmt.Errorf("eng must be non-nil")
	}

	if maxDatasets <= 0 {
		maxDatasets = 100
	}

	rows, err := eng.DatabaseHandle.Query("SELECT DatasetId, DisplayName FROM Datasets ORDER BY DatasetId LIMIT ?", maxDatasets)
	if err != nil {
		return nil, fmt.Errorf("failed to query for datasetId with error: %v", err)
	}
	defer rows.Close()

	var results []*Dataset
	for rows.Next() {
		var datasetId int64
		var displayName string
		if err := rows.Scan(&datasetId, &displayName); err != nil {
			return nil, fmt.Errorf("failed to Scan(DatasetId, DisplayName) for dataset with error: %v", err)
		}
		results = append(results, &Dataset{
			DatasetId:   datasetId,
			DisplayName: displayName,
			eng:         eng,
		})
	}
	return results, nil
}

func (d *Dataset) SetHeaders(headers bool) error {
	if d.eng == nil {
		return fmt.Errorf("eng must be non-nil")
	}
	r := "TRUE"
	if !headers {
		r = "FALSE"
	}
	stmt, err := d.eng.DatabaseHandle.Prepare("UPDATE Datasets SET HeadersSet = ? WHERE DatasetId = ?")
	if err != nil {
		return fmt.Errorf("failed to create update Datasets prepared statement with error: %v", err)
	}
	_, err = stmt.Exec(r, d.DatasetId)
	if err != nil {
		return fmt.Errorf("failed to update into Datasets table with error: %v", err)
	}
	d.HeadersSet = headers
	return nil
}
