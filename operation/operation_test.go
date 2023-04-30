package operation_test

import (
    "testing"

    "github.com/dantespe/spectacle/operation"
)

func TestNew(t *testing.T) {
    o, err := operation.New(0)
    if err != nil {
        t.Fatalf("got unexpected error from New(0): %s", err)
    }

    if o.Id != 0 {
        t.Errorf("got Id %d, wanted: 0", o.Id)
    }

    if o.Status != operation.NOT_STARTED {
        t.Errorf("got Status %d, wanted: %d", o.Status, operation.NOT_STARTED)
    }
}

func TestMarkedRunning(t *testing.T) {
    testCases := []struct {
        desc	string
        o *operation.Operation
    }{
        {
            desc: "NOT_STARTED_OP",
            o: &operation.Operation{
                Status: operation.NOT_STARTED,
            },
            
        },
        {
            desc: "RUNNING_OP",
            o: &operation.Operation{
                Status: operation.RUNNING,
            },
        },
    }
    for _, tc := range testCases {
        t.Run(tc.desc, func(t *testing.T) {
            if err := tc.o.MarkRunning(); err != nil  {
                t.Errorf("got unexpected error: %s", err)
            }
            if tc.o.Status != operation.RUNNING {
                t.Errorf("got %d, wanted: %d", tc.o.Status, operation.RUNNING)
            }
        })
    }
}

func TestMarkedRunningInvalid(t *testing.T) {
    testCases := []struct {
        desc	string
        o *operation.Operation	
    }{
        {
            desc: "COMPLETED_OP",
            o: &operation.Operation{
                Status: operation.COMPLETED,
            },
        },
        {
            desc: "FAILED_OP",
            o: &operation.Operation{
                Status: operation.FAILED,
            },
        },
    }
    for _, tc := range testCases {
        t.Run(tc.desc, func(t *testing.T) {
            if err := tc.o.MarkRunning(); err == nil {
                t.Errorf("got nil error for malformed input on MarkedRunning(): %d", tc.o.Status)
            }
        })
    }
}

func TestCompleted(t *testing.T) {
    testCases := []struct {
        desc	string
        o *operation.Operation
    }{
        {
            desc: "NOT_STARTED_OP",
            o: &operation.Operation{
                Status: operation.NOT_STARTED,
            },
        },
        {
            desc: "RUNNING_OP",
            o: &operation.Operation{
                Status: operation.RUNNING,
            },
        },
        {
            desc: "COMPLETED_OP",
            o: &operation.Operation{
                Status: operation.COMPLETED,
            },
        },
    }
    for _, tc := range testCases {
        t.Run(tc.desc, func(t *testing.T) {
            if err := tc.o.MarkCompleted(); err != nil {
                t.Fatalf("got unexpected error on MarkCompleted()")
            }
            if tc.o.Status != operation.COMPLETED {
                t.Errorf("got Status %d, wanted Status: %d", tc.o.Status, operation.COMPLETED)
            }
        })
    }
}

func TestMarkedFailed(t *testing.T) {
    testCases := []struct {
        desc	string
        o *operation.Operation
    }{
        {
            desc: "NOT_STATRTED_OP",
            o: &operation.Operation{
                Status: operation.NOT_STARTED,
            },
        },
        {
            desc: "RUNNING_OP",
            o: &operation.Operation{
                Status: operation.RUNNING,
            },
        },
        {
            desc: "COMPLETED_OP",
            o: &operation.Operation{
                Status: operation.COMPLETED,
            },
        },
        {
            desc: "FAILED_OP",
            o: &operation.Operation{
            Status: operation.FAILED,
            },
        },
    }
    for _, tc := range testCases {
        t.Run(tc.desc, func(t *testing.T) {
            msg := "Failed due to User Error"
            if err := tc.o.MarkFailed(msg); err != nil {
                t.Fatalf("Got unexpected error on MarkFailed(%s) for valid input: %s", msg, err)
            } 
        })
    }
}