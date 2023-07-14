package header

import (
	"fmt"

	"github.com/dantespe/spectacle/db"
)

type ValueType string

const (
	ValueType_RAW   ValueType = "RAW"
	ValueType_INT             = "INT"
	ValueType_FLOAT           = "FLOAT"
)

type Header struct {
	HeaderId    int64 `json:"headerId"`
	columnIndex int64
	DisplayName string `json:"displayName"`
	datasetId   int64
	valueType   ValueType
	eng         *db.Engine
}

const BucketIncrement = 1000

func (h *Header) SetColumnIndex(i int64) error {
	if h.eng == nil {
		return fmt.Errorf("cannot create a new Header with nil db.Engine")
	}
	stmt, err := h.eng.DatabaseHandle.Prepare("UPDATE Headers SET ColumnIndex = $1 WHERE HeaderId = $2")
	if err != nil {
		return fmt.Errorf("failed to create Headers prepared statement with error: %v", err)
	}
	_, err = stmt.Exec(BucketIncrement*i, h.HeaderId)
	if err != nil {
		return fmt.Errorf("failed to Update Headers table with error: %v", err)
	}
	h.columnIndex = i
	return nil
}

// New extends the dataset's headers by one and returns it.
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
	rows, err := eng.DatabaseHandle.Query("SELECT HeaderId, ValueType, DisplayName FROM Headers WHERE DatasetId = $1 ORDER BY ColumnIndex, HeaderId", datasetId)
	if err != nil {
		return nil, fmt.Errorf("failed to get headers(datasetId=%d) with error: %v", datasetId, err)
	}
	defer rows.Close()

	for rows.Next() {
		h := &Header{
			eng: eng,
		}
		if err := rows.Scan(&h.HeaderId, &h.valueType, &h.DisplayName); err != nil {
			return nil, fmt.Errorf("failed to Headers Scan with error: %v", err)
		}
		results = append(results, h)
	}
	return results, nil
}
