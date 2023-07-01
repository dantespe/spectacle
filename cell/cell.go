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
	c := &Cell{
		recordId: recordId,
		headerId: headerId,
		eng:      eng,
	}
	if err := eng.DatabaseHandle.QueryRow("INSERT INTO Cells(RecordId, HeaderId, OperationId, RawValue) VALUES($1, $2, $3, $4) RETURNING CellId", recordId, headerId, operationId, rv).Scan(&c.CellId); err != nil {
		return nil, fmt.Errorf("failed to create Cell with error: %v", err)
	}
	return c, nil
}
