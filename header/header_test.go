package header_test

import (
	"testing"

	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/header"
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

	header, err := header.New(tmp.Engine, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to create New header: %v", err)
	}

	if header.HeaderId == 0 {
		t.Errorf("got HeaderId: 0, wanted: non-zero")
	}
}

func TestGetHeaders(t *testing.T) {
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	ds, err := dataset.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create New dataset: %v", err)
	}

	for i := 0; i < 3; i++ {
		if _, err := header.New(tmp.Engine, ds.DatasetId); err != nil {
			t.Fatalf("failed to create Header: %v", err)
		}
	}

	_, err = header.GetHeaders(tmp.Engine, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to GetHeaders(%d) with err: %v", ds.DatasetId, err)
	}
}
