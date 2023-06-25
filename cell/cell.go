package cell

import (
	"fmt"

	"github.com/dantespe/spectacle/db"
)

type Cell struct {
	CellId   int64
	recordId int64
	headerId int64
	eng      *db.Engine
}

func New(eng *db.Engine, recordId int64, headerId int64, operationId int64, rv string) (*Cell, error) {
	if eng == nil {
		return nil, fmt.Errorf("cannot create a new Header with nil db.Engine")
	}
	stmt, err := eng.DatabaseHandle.Prepare("INSERT INTO Cells(RecordId, HeaderId, OperationId, RawValue) VALUES(?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("failed to build Cells PrepareStatement with error: %v", err)
	}
	res, err := stmt.Exec(recordId, headerId, operationId, rv)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into Cells table with error: %v", err)
	}
	cellId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retreive CellId with error: %v", cellId)
	}
	return &Cell{
		CellId:   cellId,
		recordId: recordId,
		headerId: headerId,
		eng:      eng,
	}, nil
}
