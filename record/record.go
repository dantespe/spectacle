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
	stmt, err := eng.DatabaseHandle.Prepare("INSERT INTO Records(DatasetId) VALUES(?)")
	if err != nil {
		return nil, fmt.Errorf("failed to build records PrepareStatement with error: %v", err)
	}
	res, err := stmt.Exec(datasetId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into Records table with error: %v", err)
	}
	recordId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retreive RecordId with error: %v", recordId)
	}
	return &Record{
		RecordId:  recordId,
		datasetId: datasetId,
		eng:       eng,
	}, nil
}
