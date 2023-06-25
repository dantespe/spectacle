package header

import (
	"fmt"

	"github.com/dantespe/spectacle/db"
)

type ValueType string

const (
	ValueType_RAW   = "RAW"
	ValueType_INT   = "INT"
	ValueType_FLOAT = "FLOAT"
)

type Header struct {
	HeaderId    int64
	datasetId   int64
	displayName string
	valueType   ValueType
	eng         *db.Engine
}

func New(eng *db.Engine, datasetId int64) (*Header, error) {
	if eng == nil {
		return nil, fmt.Errorf("cannot create a new Header with nil db.Engine")
	}
	stmt, err := eng.DatabaseHandle.Prepare("INSERT INTO Headers(DatasetId, ValueType, DisplayName) VALUES(?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("failed to build operation PrepareStatement with error: %v", err)
	}
	res, err := stmt.Exec(datasetId, ValueType_RAW, "")
	if err != nil {
		return nil, fmt.Errorf("failed to insert into Headers table with error: %v", err)
	}
	headerId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retreive HeaderId with error: %v", headerId)
	}
	return &Header{
		HeaderId:  headerId,
		datasetId: datasetId,
		valueType: ValueType_RAW,
		eng:       eng,
	}, nil
}

func GetHeaders(eng *db.Engine, datasetId int64) ([]*Header, error) {
	if eng == nil {
		return nil, fmt.Errorf("cannot GetHeaders with nil db.Engine")
	}

	var results []*Header
	rows, err := eng.DatabaseHandle.Query("SELECT HeaderId, ValueType FROM Headers WHERE DatasetId = ? ORDER BY HeaderId", datasetId)
	if err != nil {
		return nil, fmt.Errorf("failed to get headers(datasetId=%d) with error: %v", datasetId, err)
	}
	defer rows.Close()

	for rows.Next() {
		var headerId int64
		var valueType ValueType
		if err := rows.Scan(&headerId, &valueType); err != nil {
			return nil, fmt.Errorf("failed to Headers Scan with error: %v", err)
		}
		results = append(results, &Header{
			HeaderId:    headerId,
			datasetId:   datasetId,
			displayName: "",
			valueType:   valueType,
			eng:         eng,
		})
	}
	return results, nil
}
