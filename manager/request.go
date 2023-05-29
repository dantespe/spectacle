package manager

import (
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestBuilder struct{}

// CreateDatasetRequest
type CreateDatasetRequest struct {
	DisplayName string `json:"displayName"`
	HasHeaders  bool   `json:"hasHeaders"`
}

// CreateDatasetRequestBuilder from gin.Context.
func (*RequestBuilder) CreateDatasetRequestBuilder(c *gin.Context) (*CreateDatasetRequest, *CreateDatasetResponse) {
	var req CreateDatasetRequest
	if err := c.BindJSON(&req); err != nil {
		return nil, &CreateDatasetResponse{
			Message: err.Error(),
		}
	}
	return &req, nil
}

// GetDatasetRequest
type GetDatasetRequest struct {
	DatasetId uint64 `json:"datasetId"`
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
	MaxDatasets int64 `json:"maxDatasets"`
}

func (*RequestBuilder) ListDatasetsRequestBuilder(c *gin.Context) (*ListDatasetsRequest, error) {
	req := newListDatasetsRequest()
	if err := c.BindJSON(&req); err != nil {
	}

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
	// DatasetId
	DatasetId uint64 `json:"datasetId"`

	// InputFile
	InputFile io.Reader `json:"-"`
}

func (*RequestBuilder) UploadDatasetRequestBuilder(c *gin.Context) (*UploadDatasetRequest, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	// Get File from the form
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		return nil, err
	}
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	return &UploadDatasetRequest{
		DatasetId: id,
		InputFile: file,
	}, nil
}
