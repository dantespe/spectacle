// Package operation is a class that implements long running operations (LROs).
package operation

import (
	"fmt"
)

// Operation is a class that stores the state of a LRO.
type Operation struct {
	Id uint64
	Status Status
	Message string
}

// Status of the Operation.
type Status int64
const (
	NOT_STARTED Status = iota
	RUNNING 
	COMPLETED
	FAILED
)

// New returns a new Operation with the provided id.
func New(id uint64) (*Operation, error) {
	return &Operation{
		Id: id,
		Status: NOT_STARTED,
	}, nil
}

// MarkRunning sets the Status to Running. 
func (o *Operation) MarkRunning() error {
	if (o.Status == COMPLETED || o.Status == FAILED) {
		return fmt.Errorf("cannot restart a failed or completed operation")
	}

	o.Status = RUNNING
	return nil
}

// MarkCompleted sets the Status to Completed.
func (o *Operation) MarkCompleted() error {
	if (o.Status == NOT_STARTED || o.Status == RUNNING ) {
		o.Status = COMPLETED
	}
	return nil 
}

// MarkFailed sets the Status to Failed.
func (o *Operation) MarkFailed(message string) error {
	o.Status = FAILED
	o.Message = message
	return nil
}