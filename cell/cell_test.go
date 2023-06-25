package cell_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/dantespe/spectacle/cell"
	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/header"
	"github.com/dantespe/spectacle/operation"
	"github.com/dantespe/spectacle/record"
	spectesting "github.com/dantespe/spectacle/testing"
)

func CellCmp(l *cell.Cell, r *cell.Cell) bool {
	return cmp.Equal(l, r, cmpopts.IgnoreUnexported(cell.Cell{}))
}

func CellDiff(l *cell.Cell, r *cell.Cell) string {
	return cmp.Diff(l, r, cmpopts.IgnoreUnexported(cell.Cell{}))
}

func TestNew(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	ds, err := dataset.New(eng)
	if err != nil {
		t.Fatalf("failed to create New dataset: %v", err)
	}

	var records []*record.Record
	for i := 0; i < 3; i++ {
		rc, err := record.New(eng, ds.DatasetId)
		if err != nil {
			t.Fatalf("failed to create new record with error: %v", err)
		}
		records = append(records, rc)
	}

	header, err := header.New(eng, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to create Header(%d) with err: %v", ds.DatasetId, err)
	}

	op, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create a Operation with error: %v", err)
	}

	for i := 0; i < 3; i++ {
		cl, err := cell.New(eng, records[i].RecordId, header.HeaderId, op.OperationId, fmt.Sprintf("raw-value-%d", i))
		if err != nil {
			t.Fatalf("failed to create a cell with error: %v", err)
		}

		expected := &cell.Cell{CellId: int64(i + 1)}
		if !CellCmp(cl, expected) {
			t.Errorf("got diff on cell: %s", CellDiff(cl, expected))
		}
	}
}
