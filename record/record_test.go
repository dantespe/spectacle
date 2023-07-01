package record_test

import (
	"testing"

	"github.com/dantespe/spectacle/dataset"
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

	rc, err := record.New(tmp.Engine, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to create new record with error: %v", err)
	}
	if rc.RecordId == 0 {
		t.Errorf("got RecordId: 0, want: non-zero")
	}
}
