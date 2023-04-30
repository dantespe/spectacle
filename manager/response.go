package manager

import (
	"fmt"
	"net/http"

	"github.com/dantespe/spectacle/dataset"
)

type Response interface {
	JSON() map[string]interface{}
	ResponseCode() int
}

// DefaultResponse
type DefaultResponse struct {
	Response
}

func (*DefaultResponse) JSON() map[string]string {
	return make(map[string]string)
}

func (*DefaultResponse) ResponseCode() int  {
	return http.StatusOK
}

// StatusResponse
type StatusResponse struct {
	DefaultResponse

	Error error
	NumRecords int
}

func (r *StatusResponse) JSON() map[string]string {
	if r.Error != nil {
		return map[string]string{
			"status": "UNHEALTHY",
		}
	}
	return map[string]string{
		"status": "HEALTHY",
		"num_records": fmt.Sprintf("%d", r.NumRecords),
	}
}

func (r *StatusResponse) ResponseCode() int  {
	return http.StatusOK
}

// CreateDatasetResponse 
type CreateDatasetResponse struct {
	DefaultResponse

	Error error
	DatasetId uint64
}

func (r *CreateDatasetResponse) JSON() map[string]string {
	if r.Error != nil {
		return map[string]string{
			"message": "Failed to create dataset.",
		}
	}
	return map[string]string{
		"DatasetId": fmt.Sprintf("%d", r.DatasetId),
	}	
}

func (r *CreateDatasetResponse) ResponseCode() int {
	if r.Error != nil {
		return http.StatusInternalServerError
	}
	return http.StatusCreated
}

// GetDatasetResponse
type GetDatasetResponse struct {
	DefaultResponse

	Error error
	Dataset *dataset.Dataset
	DatasetId uint64
	DisplayName string
	NumRecords uint64
}

func (r *GetDatasetResponse) JSON() map[string]string {
	if r.Error != nil {
		return map[string]string {
			"message": "Dataset not found.",
		}
	}

	return map[string]string {
		"DatasetId": fmt.Sprintf("%d", r.Dataset.Id),
		"DisplayName": r.Dataset.DisplayName,
		"NumRecords": fmt.Sprintf("%d", r.Dataset.NumRecords), 
	}
}

func (r *GetDatasetResponse) ResponseCode() int {
	if r.Error != nil {
		return http.StatusNotFound
	}
	return http.StatusOK
}

// ListDatasetsResponse
type ListDatasetsResponse struct {
	DefaultResponse

	Error error
	Results []interface{}

	// Page int
	// Next string
	TotalDatasets int
}

func newListDatasetsResponse() *ListDatasetsResponse {
	return &ListDatasetsResponse{
		Results: make([]interface{}, 0),
	}
}

func (r *ListDatasetsResponse) JSON() map[string]interface{} {
	if r.Error != nil {
		return map[string]interface{} {
			"message": "Failed to list datasets.",
		}
	}

	return map[string]interface{} {
		"Results": r.Results,
		"TotalDatasets": len(r.Results),
	}
}

func (r *ListDatasetsResponse) ResponseCode() int {
	if r.Error != nil {
		return http.StatusBadRequest
	}
	return http.StatusOK
}