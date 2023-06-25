// Package operation is a class that implements long running operations (LROs).
package operation

import (
	"fmt"

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
	stmt, err := eng.DatabaseHandle.Prepare("INSERT INTO Operations(OperationStatus) VALUES(?)")
	if err != nil {
		return nil, fmt.Errorf("failed to build operation PrepareStatement with error: %v", err)
	}
	res, err := stmt.Exec(Status_NOT_STARTED)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into Operations table with error: %v", err)
	}
	operationId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retreive operationId with error: %v", operationId)
	}
	return &Operation{
		OperationId:     operationId,
		OperationStatus: Status_NOT_STARTED,
		eng:             eng,
	}, nil
}

func (o *Operation) markStatus(st Status, errMsg string) error {
	if o.eng == nil {
		return fmt.Errorf("cannot mark status when engine is nil")
	}
	stmt, err := o.eng.DatabaseHandle.Prepare("UPDATE Operations SET OperationStatus = ?, ErrorMessage = ? WHERE OperationId = ?")
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
	if o.OperationStatus == Status_SUCCESS || o.OperationStatus == Status_FAILED {
		return fmt.Errorf("cannot mark running on completed operation")
	}
	return o.markStatus(Status_RUNNING, "")
}

// MarkCompleted sets the Status to SUCCESS.
func (o *Operation) MarkSuccess() error {
	if o.OperationStatus == Status_NOT_STARTED || o.OperationStatus == Status_RUNNING {
		return o.markStatus(Status_SUCCESS, "")
	}
	return nil
}

// MarkFailed sets the Status to Failed.
func (o *Operation) MarkFailed(msg string) error {
	return o.markStatus(Status_FAILED, msg)
}
