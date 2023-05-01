package manager

import (
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
)

type RequestBuilder struct {}

// CreateDatasetRequest 
type CreateDatasetRequest struct {
	DisplayName string `json:"displayName"`
}

// CreateDatasetRequestBuilder from gin.Context.
func (*RequestBuilder) CreateDatasetRequestBuilder(c *gin.Context) (*CreateDatasetRequest, error) {
	var req CreateDatasetRequest
	if err := c.BindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// GetDatasetRequest
type GetDatasetRequest struct {
	DatasetId uint64
}

func (*RequestBuilder) GetDatasetRequestBuilder(c *gin.Context) (*GetDatasetRequest, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return &GetDatasetRequest{
		DatasetId: id,
	}, nil
}

// ListDatasets
type ListDatasetsRequest struct {
	// Defaults to 1000
	// UNUSED 
	// TODO: Turn this on
	MaxDatasets int64
}

func (*RequestBuilder) ListDatasetsRequestBuilder(c *gin.Context) (*ListDatasetsRequest, error) {
	req := newListDatasetsRequest()

	if md := c.Param("max_datasets"); md != "" {
		m, err := strconv.ParseInt(md, 10, 32)
		if err != nil {
			return nil, err
		}
		if m < 0 {
			return nil, fmt.Errorf("max_datasets must be a postive number")
		}
		req.MaxDatasets = m
	}

	return req, nil
}

func newListDatasetsRequest() *ListDatasetsRequest {
	return &ListDatasetsRequest{
		MaxDatasets: 1000,
	}
}

// UploadDatasetRequest
type UploadDatasetRequest struct {
	// Default csv 
	ConnectorType string
	
	// Default True
	// TODO
	HasHeaders bool
}

const (
	ConnectorType_CSV string = "csv"
)

func (*RequestBuilder) UploadDatasetRequestBuilder(*gin.Context) (*UploadDatasetRequest, error) {
	// TODO: Read the params from c
	return &UploadDatasetRequest{
		ConnectorType: ConnectorType_CSV,
		HasHeaders: true,
	}, nil
}