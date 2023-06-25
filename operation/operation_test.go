package operation_test

import (
	"os"
	"testing"

	"github.com/dantespe/spectacle/operation"
	spectesting "github.com/dantespe/spectacle/testing"
)

func TestNew(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	o, err := operation.New(eng)
	if err != nil {
		t.Fatalf("got unexpected error from New(): %s", err)
	}
	if o.OperationId != 1 {
		t.Errorf("got Id %d, wanted: 1", o.OperationId)
	}

	if o.OperationStatus != operation.Status_NOT_STARTED {
		t.Errorf("got Status %s, wanted: %s", o.OperationStatus, operation.Status_NOT_STARTED)
	}
}

func TestMarkedRunning(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	// Not Started Op
	op, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create NOT_STARTED op with err: %v", err)
	}
	if op.OperationStatus != operation.Status_NOT_STARTED {
		t.Fatalf("got OperationStatus: %s, want: %s", op.OperationStatus, operation.Status_NOT_STARTED)
	}

	// Running Op
	op2, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create NOT_STARTED op2 with err: %v", err)
	}
	op2.OperationStatus = operation.Status_RUNNING

	testCases := []struct {
		desc string
		o    *operation.Operation
	}{
		{
			desc: "not_started_op",
			o:    op,
		},
		{
			desc: "running_op",
			o:    op2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if err := tc.o.MarkRunning(); err != nil {
				t.Errorf("got unexpected error for MarkRunning(): %v", err)
			}
			if tc.o.OperationStatus != operation.Status_RUNNING {
				t.Errorf("got %s, wanted: %s", tc.o.OperationStatus, operation.Status_RUNNING)
			}
		})
	}
}

func TestMarkedRunningInvalid(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	op, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create operation with err: %v", err)
	}
	if err := op.MarkSuccess(); err != nil {
		t.Fatalf("failed to set status: %v", err)
	}

	op2, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create operation with err: %v", err)
	}
	op2.MarkFailed("failed to due to error injection")

	testCases := []struct {
		desc string
		o    *operation.Operation
	}{
		{
			desc: "SUCCESS_OP",
			o:    op,
		},
		{
			desc: "FAILED_OP",
			o:    op2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if err := tc.o.MarkRunning(); err == nil {
				t.Errorf("got nil error for malformed input on MarkedRunning(): %s", tc.o.OperationStatus)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	nsOp, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}

	runOp, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}
	if err := runOp.MarkRunning(); err != nil {
		t.Fatalf("failed to set status to running: %v", err)
	}

	succOp, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}
	if err := succOp.MarkSuccess(); err != nil {
		t.Fatalf("failed to set status to success: %v", err)
	}

	testCases := []struct {
		desc string
		o    *operation.Operation
	}{
		{
			desc: "NOT_STARTED_OP",
			o:    nsOp,
		},
		{
			desc: "RUNNING_OP",
			o:    runOp,
		},
		{
			desc: "COMPLETED_OP",
			o: &operation.Operation{
				OperationStatus: operation.Status_SUCCESS,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if err := tc.o.MarkSuccess(); err != nil {
				t.Fatalf("got unexpected error on MarkCompleted(): %v", err)
			}
			if tc.o.OperationStatus != operation.Status_SUCCESS {
				t.Errorf("got Status %s, wanted Status: %s", tc.o.OperationStatus, operation.Status_SUCCESS)
			}
		})
	}
}

func TestMarkedFailed(t *testing.T) {
	eng, fname, err := spectesting.CreateTempSQLiteEngine()
	if err != nil {
		t.Fatalf("failed to create database engine: %v", err)
	}
	defer os.Remove(fname)

	nsOp, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}

	runOp, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}
	if err := runOp.MarkRunning(); err != nil {
		t.Fatalf("failed to set status to running: %v", err)
	}

	succOp, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}
	if err := succOp.MarkSuccess(); err != nil {
		t.Fatalf("failed to set status to success: %v", err)
	}

	fOp, err := operation.New(eng)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}
	if err := fOp.MarkFailed("failed due to user injection"); err != nil {
		t.Fatalf("failed to set status to success: %v", err)
	}

	testCases := []struct {
		desc string
		o    *operation.Operation
	}{
		{
			desc: "NOT_STATRTED_OP",
			o:    nsOp,
		},
		{
			desc: "RUNNING_OP",
			o:    runOp,
		},
		{
			desc: "SUCCESS_OP",
			o:    succOp,
		},
		{
			desc: "FAILED_OP",
			o:    fOp,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			msg := "Failed due to User Error"
			if err := tc.o.MarkFailed(msg); err != nil {
				t.Fatalf("Got unexpected error on MarkFailed(%s) for valid input: %s", msg, err)
			}
			if tc.o.OperationStatus != operation.Status_FAILED {
				t.Errorf("got %s, want: %s", tc.o.OperationStatus, operation.Status_FAILED)
			}
		})
	}
}
