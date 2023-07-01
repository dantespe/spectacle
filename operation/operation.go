// Package operation is a class that implements long running operations (LROs).
package operation

import (
	"fmt"
	"sync"

	"github.com/dantespe/spectacle/db"
)

// Status of the Operation.
type Status string

const (
	Status_NOT_STARTED Status = "NOT_STARTED"
	Status_RUNNING     Status = "RUNNING"
	Status_SUCCESS     Status = "SUCCESS"
	Status_FAILED      Status = "FAILED"
)

// Operation is a class that stores the state of a LRO.
type Operation struct {
	mu              sync.Mutex
	OperationId     int64
	OperationStatus Status
	ErrorMessage    string
	eng             *db.Engine
}

// New creates a new Operation and saves it to the database.
func New(eng *db.Engine) (*Operation, error) {
	if eng == nil {
		return nil, fmt.Errorf("cannot create new operation with nil db.Engine")
	}

	op := &Operation{
		OperationStatus: Status_NOT_STARTED,
		eng:             eng,
	}

	if err := eng.DatabaseHandle.QueryRow("INSERT INTO Operations(OperationStatus) VALUES($1) RETURNING OperationId", Status_NOT_STARTED).Scan(&op.OperationId); err != nil {
		return nil, fmt.Errorf("failed to create operation with error: %v", err)
	}
	return op, nil
}

func (o *Operation) markStatus(st Status, errMsg string) error {
	if o.eng == nil {
		return fmt.Errorf("cannot mark status when engine is nil")
	}
	stmt, err := o.eng.DatabaseHandle.Prepare("UPDATE Operations SET OperationStatus = $1, ErrorMessage = $2 WHERE OperationId = $3")
	if err != nil {
		return fmt.Errorf("failed to build operation PrepareStatement with error: %v", err)
	}
	_, err = stmt.Exec(st, errMsg, o.OperationId)
	if err != nil {
		return fmt.Errorf("failed to update operations table with error: %v", err)
	}
	o.OperationStatus = st
	o.ErrorMessage = errMsg
	return nil
}

// MarkRunning sets the OperationStatus to RUNNING.
func (o *Operation) MarkRunning() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.OperationStatus == Status_SUCCESS || o.OperationStatus == Status_FAILED {
		return fmt.Errorf("cannot mark running on completed operation")
	}
	return o.markStatus(Status_RUNNING, "")
}

// MarkCompleted sets the Status to SUCCESS.
func (o *Operation) MarkSuccess() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.OperationStatus == Status_NOT_STARTED || o.OperationStatus == Status_RUNNING {
		return o.markStatus(Status_SUCCESS, "")
	}
	return nil
}

// MarkFailed sets the Status to Failed.
func (o *Operation) MarkFailed(msg string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.markStatus(Status_FAILED, msg)
}

func (o *Operation) Complete() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.OperationStatus == Status_FAILED || o.OperationStatus == Status_SUCCESS
}
