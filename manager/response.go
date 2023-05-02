package manager

import (
    "net/http"

    "github.com/dantespe/spectacle/dataset"
)

type JSONResponse map[string]interface{}

type Response interface {
    JSON() JSONResponse
    ResponseCode() int
}

// DefaultResponse
type DefaultResponse struct {
    Response
}

func (*DefaultResponse) JSON() JSONResponse {
    return make(JSONResponse)
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

func (r *StatusResponse) JSON() JSONResponse {
    if r.Error != nil {
        return JSONResponse {
            "status": "UNHEALTHY",
        }
    }
    return JSONResponse {
        "status": "HEALTHY",
        "num_records": r.NumRecords,
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
    DatasetUrl string
    DisplayName string
}

func (r *CreateDatasetResponse) JSON() JSONResponse {
    if r.Error != nil {
        return map[string]interface{}{
            "message": "Failed to create dataset.",
        }
    }
    return map[string]interface{}{
        "DatasetUrl": r.DatasetUrl,
        "DatasetId": r.DatasetId,
        "DisplayName": r.DisplayName,
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

func (r *GetDatasetResponse) JSON() JSONResponse {
    if r.Error != nil {
        return JSONResponse {
            "message": "Dataset not found.",
        }
    }

    return JSONResponse {
        "DatasetId": r.Dataset.Id,
        "DisplayName": r.Dataset.DisplayName,
        "NumRecords": r.Dataset.NumRecords,
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

func (r *ListDatasetsResponse) JSON() JSONResponse {
    if r.Error != nil {
        return JSONResponse {
            "message": "Failed to list datasets.",
        }
    }

    return JSONResponse {
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