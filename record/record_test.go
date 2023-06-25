package record_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/record"
	spectesting "github.com/dantespe/spectacle/testing"
)

func RecordCmp(l *record.Record, r *record.Record) bool {
	return cmp.Equal(l, r, cmpopts.IgnoreUnexported(record.Record{}))
}

func RecordDiff(l *record.Record, r *record.Record) string {
	return cmp.Diff(l, r, cmpopts.IgnoreUnexported(record.Record{}))
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

	rc, err := record.New(eng, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to create new record with error: %v", err)
	}

	expected := &record.Record{RecordId: 1}

	if !RecordCmp(rc, expected) {
		t.Errorf("got unexpected diff for record.New: %s", RecordDiff(rc, expected))
	}
}
