package record

import (
	"fmt"

	"github.com/dantespe/spectacle/db"
)

type Record struct {
	RecordId  int64
	datasetId int64
	eng       *db.Engine
}

func New(eng *db.Engine, datasetId int64) (*Record, error) {
	if eng == nil {
		return nil, fmt.Errorf("cannot create a new Header with nil db.Engine")
	}
	r := &Record{
		datasetId: datasetId,
		eng:       eng,
	}
	if err := eng.DatabaseHandle.QueryRow("INSERT INTO Records(DatasetId) VALUES($1) RETURNING RecordId", datasetId).Scan(&r.RecordId); err != nil {
		return nil, fmt.Errorf("failed to create Record with error: %v", err)
	}
	return r, nil
}
