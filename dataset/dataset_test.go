// Tests for github.com/dantespe/spectacle/dataset.
package dataset_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/dantespe/spectacle/dataset"
)

func TestDefault(t *testing.T) {
	ds, err := dataset.Default()
	if err != nil {
		t.Fatalf("Got unexpected err for Default(): %s", err)
	}
	if ds.Id != 0 {
		t.Errorf("Got Id: %d, want Id: 0", ds.Id)
	}
	if ds.NumRecords != 0 {
		t.Errorf("Got NumRecords: %d, want: 0", ds.NumRecords)
	}
}

func TestNewWithId(t *testing.T) {
	testCases := []struct {
		desc	string
		id uint64
	}{
		{
			desc: "default",
			id: 0,
		},
		{
			desc: "specified_value",
			id: 56985,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ds, err := dataset.NewWithId(tc.id)
			if err != nil {
				t.Fatalf("Got unexpected error on valid input: %s", err)
			}
			if ds.Id != tc.id {
				t.Errorf("Got Id: %d, wanted Id: %d", ds.Id, tc.id)
			}
		})
	}
}

func TestSummary(t *testing.T) {
	testCases := []struct {
		desc	string
		ds *dataset.Dataset
		expectedSummary map[string]interface{}
	}{
		{
			desc: "default",
			ds: &dataset.Dataset{
				Id: 123,
				DisplayName: "my-dataset",
				NumRecords: 10000,
			},
			expectedSummary: map[string]interface{}{
				"datasetId": uint64(123),
				"displayName": "my-dataset",
				"numRecords": uint64(10000),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if diff := deep.Equal(tc.ds.Summary(), tc.expectedSummary); diff != nil {
				t.Error(diff)
			}
		})
	}
}