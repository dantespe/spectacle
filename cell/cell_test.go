package cell_test

import (
	"fmt"
	"testing"

	"github.com/dantespe/spectacle/cell"
	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/header"
	"github.com/dantespe/spectacle/operation"
	"github.com/dantespe/spectacle/record"
	spectesting "github.com/dantespe/spectacle/testing"
)

func TestNew(t *testing.T) {
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	ds, err := dataset.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create New dataset: %v", err)
	}

	var records []*record.Record
	for i := 0; i < 3; i++ {
		rc, err := record.New(tmp.Engine, ds.DatasetId)
		if err != nil {
			t.Fatalf("failed to create new record with error: %v", err)
		}
		records = append(records, rc)
	}

	header, err := header.New(tmp.Engine, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to create Header(%d) with err: %v", ds.DatasetId, err)
	}

	op, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create a Operation with error: %v", err)
	}

	for i := 0; i < 3; i++ {
		_, err := cell.New(tmp.Engine, records[i].RecordId, header.HeaderId, op.OperationId, fmt.Sprintf("raw-value-%d", i))
		if err != nil {
			t.Fatalf("failed to create a cell with error: %v", err)
		}
	}
}
