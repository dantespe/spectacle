package header_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/header"
	spectesting "github.com/dantespe/spectacle/testing"
)

func HeaderCmp(l *header.Header, r *header.Header) bool {
	return cmp.Equal(l, r, cmpopts.IgnoreUnexported(header.Header{}))
}

func HeaderDiff(l *header.Header, r *header.Header) string {
	return cmp.Diff(l, r, cmpopts.IgnoreUnexported(header.Header{}))
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

	header, err := header.New(eng, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to create New header: %v", err)
	}

	if header.HeaderId != 1 {
		t.Errorf("Got HeaderId: %d, want: 1", header.HeaderId)
	}
}

func TestGetHeaders(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	ds, err := dataset.New(eng)
	if err != nil {
		t.Fatalf("failed to create New dataset: %v", err)
	}

	for i := 0; i < 3; i++ {
		if _, err := header.New(eng, ds.DatasetId); err != nil {
			t.Fatalf("failed to create Header: %v", err)
		}
	}

	headers, err := header.GetHeaders(eng, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to GetHeaders(%d) with err: %v", ds.DatasetId, err)
	}

	expected := []*header.Header{
		&header.Header{
			HeaderId: 1,
		},
		&header.Header{
			HeaderId: 2,
		},
		&header.Header{
			HeaderId: 3,
		},
	}

	for i := 0; i < 3; i++ {
		if !HeaderCmp(headers[i], expected[i]) {
			t.Fatalf("got unexpected diff: %s", HeaderDiff(headers[i], expected[i]))
		}
	}
}
