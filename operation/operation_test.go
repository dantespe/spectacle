package operation_test

import (
	"testing"

	"github.com/dantespe/spectacle/operation"
	spectesting "github.com/dantespe/spectacle/testing"
)

func TestNew(t *testing.T) {
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	o, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("got unexpected error from New(): %s", err)
	}
	if o.OperationId == 0 {
		t.Errorf("got id: 0, wanted: non-zero")
	}

	if o.OperationStatus != operation.Status_NOT_STARTED {
		t.Errorf("got Status %s, wanted: %s", o.OperationStatus, operation.Status_NOT_STARTED)
	}
}

func TestMarkedRunning(t *testing.T) {
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	// Not Started Op
	op, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create NOT_STARTED op with err: %v", err)
	}
	if op.OperationStatus != operation.Status_NOT_STARTED {
		t.Fatalf("got OperationStatus: %s, want: %s", op.OperationStatus, operation.Status_NOT_STARTED)
	}

	// Running Op
	op2, err := operation.New(tmp.Engine)
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
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	op, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create operation with err: %v", err)
	}
	if err := op.MarkSuccess(); err != nil {
		t.Fatalf("failed to set status: %v", err)
	}

	op2, err := operation.New(tmp.Engine)
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
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	nsOp, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}

	runOp, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}
	if err := runOp.MarkRunning(); err != nil {
		t.Fatalf("failed to set status to running: %v", err)
	}

	succOp, err := operation.New(tmp.Engine)
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
			desc: "COMPLETE_OP",
			o:    succOp,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if err := tc.o.MarkSuccess(); err != nil {
				t.Fatalf("got unexpected error on MarkSuccess(): %v", err)
			}
			if tc.o.OperationStatus != operation.Status_SUCCESS {
				t.Errorf("got Status %s, wanted Status: %s", tc.o.OperationStatus, operation.Status_SUCCESS)
			}
		})
	}
}

func TestMarkedFailed(t *testing.T) {
	tmp, err := spectesting.NewTempPostgres()
	if err != nil {
		t.Fatalf("failed to create temp postgres database with err: %v", err)
	}
	defer tmp.Close()

	nsOp, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}

	runOp, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}
	if err := runOp.MarkRunning(); err != nil {
		t.Fatalf("failed to set status to running: %v", err)
	}

	succOp, err := operation.New(tmp.Engine)
	if err != nil {
		t.Fatalf("failed to create new operation: %v", err)
	}
	if err := succOp.MarkSuccess(); err != nil {
		t.Fatalf("failed to set status to success: %v", err)
	}

	fOp, err := operation.New(tmp.Engine)
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
