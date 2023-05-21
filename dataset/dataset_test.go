// Tests for github.com/dantespe/spectacle/dataset.
package dataset_test

import (
    "testing"

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

func TestNewWithOpts(t *testing.T) {
    testCases := []struct {
        desc	string
        opts []dataset.Option
        expectedDataset *dataset.Dataset
    }{
        {
            desc: "WithId",
            opts: []dataset.Option{
                dataset.WithId(3),
            },
            expectedDataset: &dataset.Dataset{
                Id: 3,
                HasHeaders: true,
            },
        },
        {
            desc: "WithIdAndDisplayName",
            opts: []dataset.Option{
                dataset.WithId(123),
                dataset.WithDisplayName("my-display-name"),
            },
            expectedDataset: &dataset.Dataset{
                Id: 123,
                DisplayName: "my-display-name",
                HasHeaders: true,
            },
        },
        {
            desc: "WithoutHasHeaders",
            opts: []dataset.Option{
                dataset.WithId(1),
                dataset.WithHasHeaders(false),
            },
            expectedDataset: &dataset.Dataset{
                Id: 1,
                HasHeaders: false,
            },
        },
    }
    for _, tc := range testCases {
        t.Run(tc.desc, func(t *testing.T) {
            ds, err := dataset.New(tc.opts...)
            if err != nil {
                t.Fatalf("Got unexpected err for New(opts): %s", err)
            }
            if equal, diff := ds.Equal(tc.expectedDataset); !equal {
                t.Errorf("Got unexpected diff from expected Dataset: %s", diff)
            }
        })
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

func TestCopy(t *testing.T) {
    testCases := []struct {
        desc	string
        ds *dataset.Dataset
        id uint64
        displayName string
        numRecords uint64
    }{
        {
            desc: "default",
            ds: &dataset.Dataset{
                Id: 123,
                DisplayName: "my-dataset",
                NumRecords: 10000,
            },
            id: 123,
            displayName: "my-dataset",
            numRecords: 10000,
        },
    }
    for _, tc := range testCases {
        t.Run(tc.desc, func(t *testing.T) {
            cpy := tc.ds.Copy()
            if cpy.Id != tc.ds.Id {
                t.Errorf("Got Id: %d, want: %d", cpy.Id, tc.ds.Id)
            }
            if cpy.DisplayName != tc.ds.DisplayName {
                t.Errorf("Got DisplayName: %s, want: %s", cpy.DisplayName, tc.ds.DisplayName)
            }
            if cpy.NumRecords != tc.ds.NumRecords {
                t.Errorf("Num Records: %d, want: %d", cpy.NumRecords, tc.ds.NumRecords)
            }
        })
    }
}

func TestSetUntitledDisplayName(t *testing.T) {
    testCases := []struct {
        desc	string
        displayName string
        uId int
        expectedDisplayName string 
        expectedUid int
    }{
        {
            desc: "provided",
            displayName: "test-123",
            uId: 4,
            expectedDisplayName: "test-123",
            expectedUid: 4,
        },
        {
            desc: "untitled",
            displayName: "",
            uId: 8,
            expectedDisplayName: "untitled-8",
            expectedUid: 9,
        },
    }
    for _, tc := range testCases {
        t.Run(tc.desc, func(t *testing.T) {
            ds, err := dataset.New(
                dataset.WithDisplayName(tc.displayName),
            )
            if err != nil {
                t.Fatalf("failed to create dataset")
            }
            uid := ds.SetUntitledDisplayName(tc.uId)
            if ds.DisplayName != tc.expectedDisplayName {
                t.Errorf("Got DisplayName: %s, want: %s", ds.DisplayName, tc.expectedDisplayName)
            }
            if uid != tc.expectedUid {
                t.Errorf("Got uId: %d, want: %d", uid, tc.expectedUid)
            }
        })
    }
}