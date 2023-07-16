package manager

import (
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestBuilder struct{}

// CreateDatasetRequest
type CreateDatasetRequest struct {
	DisplayName string `json:"displayName"`
}

// CreateDatasetRequestBuilder from gin.Context.
func (*RequestBuilder) CreateDatasetRequestBuilder(c *gin.Context) (*CreateDatasetRequest, *CreateDatasetResponse) {
	var req CreateDatasetRequest
	c.ShouldBindJSON(&req)
	return &req, nil
}

// GetDatasetRequest
type GetDatasetRequest struct {
	DatasetId int64 `json:"datasetId"`
}

func (*RequestBuilder) GetDatasetRequestBuilder(c *gin.Context) (*GetDatasetRequest, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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
	c.ShouldBindJSON(&req)
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
	DatasetId int64 `json:"datasetId"`

	// HasHeaders
	HasHeaders bool `json:"hasHeaders"`

	// InputFile
	InputFile io.Reader `json:"-"`
}

func (*RequestBuilder) UploadDatasetRequestBuilder(c *gin.Context) (*UploadDatasetRequest, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

type GetHeadersRequest struct {
	// DatasetId
	DatasetId int64 `json:"datasetId"`
}

func (*RequestBuilder) GetHeadersRequestBuilder(c *gin.Context) (*GetHeadersRequest, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return &GetHeadersRequest{DatasetId: id}, nil
}

type DataRequest struct {
	DatasetId    int64   `json:"datasetId"`
	Headers      []int64 `json:"headers"`
	LastRecordId int64   `json:"recordid"`
	MaxResults   int64   `json:"maxresults"`
}

func (*RequestBuilder) DataRequestBuilder(c *gin.Context) (*DataRequest, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	headers := make([]int64, 0)
	if c.Query("headers") != "" {
		for _, h := range strings.Split(c.Query("headers"), ",") {
			id, err := strconv.ParseInt(h, 10, 64)
			if err != nil {
				return nil, err
			}
			headers = append(headers, id)
		}
	}

	lastRecordId := int64(-1)
	if c.Query("recordid") != "" {
		var err error
		lastRecordId, err = strconv.ParseInt(c.Query("recordid"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	maxResults := int64(100)
	if c.Query("maxresults") != "" {
		var err error
		maxResults, err = strconv.ParseInt(c.Query("maxresults"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	resp := &DataRequest{
		DatasetId:    id,
		Headers:      headers,
		LastRecordId: lastRecordId,
		MaxResults:   maxResults,
	}
	return resp, nil
}
