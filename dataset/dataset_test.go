// Tests for github.com/dantespe/spectacle/dataset.
package dataset_test

import (
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
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

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
			ds, err := dataset.New(tmp.Engine, tc.opts...)
			if err != nil {
				t.Fatalf("failed to create New dataset: %v", err)
			}
			ds2, err := dataset.GetDatasetFromId(tmp.Engine, ds.DatasetId)
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
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	before, err := dataset.TotalDatasets(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to get total datasets with err: %v", err)
	}

	for i := 0; i < 5; i++ {
		_, err := dataset.New(tmp.Engine)
		if err != nil {
			t.Fatalf("failed to create dataset with err: %v", err)
		}
	}

	after, err := dataset.TotalDatasets(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to get total datasets with err: %v", err)
	}
	if after < before+5 {
		t.Errorf("got: %d, wanted at least: %d", after, before+5)
	}
}

func TestGetDatasets(t *testing.T) {
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	for i := 0; i < 3; i++ {
		_, err := dataset.New(tmp.Engine)
		if err != nil {
			t.Fatalf("failed to create dataset with err: %v", err)
		}
	}

	ds, err := dataset.GetDatasets(tmp.Engine, 3)
	if err != nil {
		t.Fatalf("got error for GetDatasets(3): %v", err)
	}

	if len(ds) != 3 {
		t.Errorf("got len: %d, want: 3", len(ds))
	}
}

func TestSetDatasets(t *testing.T) {
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	ds, err := dataset.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create dataset with err: %v", err)
	}

	if err := ds.SetHeaders(true); err != nil {
		t.Errorf("failed to SetHeaders with err: %v", err)
	}

	ds2, err := dataset.GetDatasetFromId(tmp.Engine, ds.DatasetId)
	if err != nil {
		t.Fatalf("failed to GetDatasetFromId with err: %v", err)
	}
	if !DatasetCmp(ds, ds2) {
		t.Errorf("Got diff: %s, want: ''", DatasetDiff(ds, ds2))
	}
}

func TestUpdateNumRecords(t *testing.T) {
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	ds, err := dataset.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create dataset with err: %v", err)
	}

	if err = ds.UpdateNumRecords(); err != nil {
		t.Fatalf("got unexpected err on UpdateNumRecords: %v", err)
	}
	if ds.NumRecords != 0 {
		t.Errorf("got NumRecords: %d, want: 0", ds.NumRecords)
	}
}
