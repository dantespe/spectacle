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

	h := &Header{
		datasetId: datasetId,
		valueType: ValueType_RAW,
		eng:       eng,
	}

	if err := eng.DatabaseHandle.QueryRow("INSERT INTO Headers(DatasetId, ValueType, DisplayName) VALUES($1, $2, $3) RETURNING HeaderId", datasetId, ValueType_RAW, "").Scan(&h.HeaderId); err != nil {
		return nil, fmt.Errorf("failed to build operation PrepareStatement with error: %v", err)
	}
	return h, nil
}

func GetHeaders(eng *db.Engine, datasetId int64) ([]*Header, error) {
	if eng == nil {
		return nil, fmt.Errorf("cannot GetHeaders with nil db.Engine")
	}

	var results []*Header
	rows, err := eng.DatabaseHandle.Query("SELECT HeaderId, ValueType FROM Headers WHERE DatasetId = $1 ORDER BY HeaderId", datasetId)
	if err != nil {
		return nil, fmt.Errorf("failed to get headers(datasetId=%d) with error: %v", datasetId, err)
	}
	defer rows.Close()

	for rows.Next() {
		h := &Header{
			eng: eng,
		}
		if err := rows.Scan(&h.HeaderId, &h.valueType); err != nil {
			return nil, fmt.Errorf("failed to Headers Scan with error: %v", err)
		}
		results = append(results, h)
	}
	return results, nil
}
