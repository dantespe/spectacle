// Tests for github.com/dantespe/spectacle/dataset.
package dataset_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/dantespe/spectacle/dataset"
	spectesting "github.com/dantespe/spectacle/testing"
)

func DatasetCmp(l *dataset.Dataset, r *dataset.Dataset) bool {
	return cmp.Equal(l, r, cmpopts.IgnoreUnexported(dataset.Dataset{}))
}

func DatasetDiff(l *dataset.Dataset, r *dataset.Dataset) string {
	return cmp.Diff(l, r, cmpopts.IgnoreUnexported(dataset.Dataset{}))
}

func TestNew(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	testCases := []struct {
		desc string
		opts []dataset.Option
	}{
		{
			desc: "no_opts",
			opts: []dataset.Option{},
		},
		{
			desc: "with_display_name",
			opts: []dataset.Option{
				dataset.WithDisplayName("super-dataset"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ds, err := dataset.New(eng, tc.opts...)
			if err != nil {
				t.Fatalf("failed to create New dataset: %v", err)
			}
			ds2, err := dataset.GetDatasetFromId(eng, ds.DatasetId)
			if err != nil {
				t.Fatalf("failed to retrieve dataset: %v", err)
			}
			if !DatasetCmp(ds, ds2) {
				t.Errorf("Got diff: %s, want: ''", DatasetDiff(ds, ds2))
			}
		})
	}
}

func TestTotalDatasets(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	for i := 0; i < 5; i++ {
		_, err := dataset.New(eng)
		if err != nil {
			t.Fatalf("failed to create dataset with err: %v", err)
		}
	}

	td, err := dataset.TotalDatasets(eng)
	if err != nil {
		t.Errorf("got: %d, want: 5", td)
	}
}

func TestGetDatasets(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	for i := 0; i < 3; i++ {
		_, err := dataset.New(eng)
		if err != nil {
			t.Fatalf("failed to create dataset with err: %v", err)
		}
	}

	expected := []*dataset.Dataset{
		&dataset.Dataset{
			DatasetId:   1,
			DisplayName: "untitled-1",
		},
		&dataset.Dataset{
			DatasetId:   2,
			DisplayName: "untitled-2",
		},
		&dataset.Dataset{
			DatasetId:   3,
			DisplayName: "untitled-3",
		},
	}

	result, err := dataset.GetDatasets(eng, 10)
	if err != nil {
		t.Fatalf("failed to GetDatasets() with err: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("Got len(result): %d, with: %d", len(result), len(expected))
	}

	for i := 0; i < 3; i++ {
		if !DatasetCmp(result[i], expected[i]) {
			t.Errorf("Got diff: %s, want: ''", DatasetDiff(result[i], expected[i]))
		}
	}
}

func TestSetDatasets(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	ds, err := dataset.New(eng)
	if err != nil {
		t.Fatalf("failed to create dataset with err: %v", err)
	}

	if err := ds.SetHeaders(true); err != nil {
		t.Errorf("failed to SetHeaders with err: %v", err)
	}

	ds2, err := dataset.GetDatasetFromId(eng, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to GetDatasetFromId with err: %v", err)
	}
	if !DatasetCmp(ds, ds2) {
		t.Errorf("Got diff: %s, want: ''", DatasetDiff(ds, ds2))
	}
}
