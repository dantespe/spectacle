package manager

import (
	"github.com/dantespe/spectacle/dataset"
	"github.com/dantespe/spectacle/header"
)

// StatusResponse
type StatusResponse struct {
	Message     string `json:"-"`
	NumRecords  int    `json:"numRecords"`
	NumDatasets int    `json:"numDatasets"`
	Status      string `json:"status"`
	Code        int    `json:"code"`
}

// CreateDatasetResponse
type CreateDatasetResponse struct {
	Message     string `json:"error,omitempty"`
	DatasetId   int64  `json:"datasetId,omitempty"`
	DatasetUrl  string `json:"datasetUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Code        int    `json:"code"`
}

// GetDatasetResponse
type GetDatasetResponse struct {
	Message string           `json:"error,omitempty"`
	Dataset *dataset.Dataset `json:"dataset,omitempty"`
	Code    int              `json:"code"`
}

// ListDatasetsResponse
type ListDatasetsResponse struct {
	Results       []*dataset.Dataset `json:"results"`
	TotalDatasets int64              `json:"totalDatasets"`
	Message       string             `json:"error,omitempty"`
	Code          int                `json:"code"`
}

// UploadDatasetResponse
type UploadDatasetResponse struct {
	OperationUrl string `json:"operation,omitempty"`
	Message      string `json:"error,omitempty"`
	Code         int    `json:"code"`
}

type GetHeadersResponse struct {
	Headers []*header.Header `json:"results"`
	Message string           `json:"error,omitempty"`
	Code    int              `json:"code"`
}

type ResultSet struct {
	HeaderId    int64    `json:"headerId"`
	DisplayName string   `json:"displayName"`
	Data        []string `json:"rows"`
}

type DataResponse struct {
	Results []*ResultSet `json:"results"`
	Next    string       `json:"next,omitempty"`
	Message string       `json:"error,omitempty"`
	Code    int          `json:"code"`
}
