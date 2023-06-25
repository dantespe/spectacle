package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dantespe/spectacle/handler"
	"github.com/dantespe/spectacle/manager"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func GetRouter() *gin.Engine {
	router := gin.Default()
	rg := router.Group("rest")
	handler.AddRestHandlerRoutes(rg)
	return router
}

func BuffFromRequest(t *testing.T, req interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		t.Fatalf("failed to encode request (%v) to json: %v", req, err)
	}
	return &buf
}

func TestStatus(t *testing.T) {
	req, _ := http.NewRequest("GET", "/rest/status", nil)

	router := GetRouter()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK, "response code")
}

func TestCreateDataset(t *testing.T) {
	tests := []struct {
		desc          string
		req           *manager.CreateDatasetRequest
		exDisplayName string
	}{
		{
			desc:          "untitled",
			req:           &manager.CreateDatasetRequest{},
			exDisplayName: "untitled-1",
		},
		{
			desc:          "untitled_again",
			req:           &manager.CreateDatasetRequest{},
			exDisplayName: "untitled-2",
		},
		{
			desc: "with_display_name",
			req: &manager.CreateDatasetRequest{
				DisplayName: "my-display-name",
			},
			exDisplayName: "my-display-name",
		},
		{
			desc: "duplicate_display_name",
			req: &manager.CreateDatasetRequest{
				DisplayName: "my-display-name",
			},
			exDisplayName: "my-display-name",
		},
	}

	router := GetRouter()
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/rest/dataset", BuffFromRequest(t, tc.req))
			if err != nil {
				t.Fatalf("failed to build request with err: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var resp manager.CreateDatasetResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal request with err: %v", err)
			}
			assert.Equal(t, w.Code, http.StatusCreated, "response code")
			assert.Equal(t, w.Code, resp.Code, "Code")
			assert.Empty(t, resp.Message) // only populated on errors
			assert.Equal(t, resp.DisplayName, tc.exDisplayName)
			assert.NotEmpty(t, resp.DatasetId)
			assert.Contains(t, resp.DatasetUrl, fmt.Sprintf("%d", resp.DatasetId))
		})
	}
}

func TestGetDataset(t *testing.T) {
	// Create a dataset
	router := GetRouter()
	req, err := http.NewRequest("POST", "/rest/dataset", BuffFromRequest(t, &manager.CreateDatasetRequest{}))
	if err != nil {
		t.Fatalf("failed to build request with err: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp manager.CreateDatasetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal request with err: %v", err)
	}
	assert.Equal(t, w.Code, http.StatusCreated, "response code")

	// Get a dataset
	req, err = http.NewRequest("GET", fmt.Sprintf("/rest/dataset/%d", resp.DatasetId), nil)
	if err != nil {
		t.Fatalf("failed to build http request with err: %v", err)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp2 manager.GetDatasetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp2); err != nil {
		t.Fatalf("failed to unmarshal json with err: %v", err)
	}

	assert.Equal(t, w.Code, http.StatusOK, "response code")
	assert.Equal(t, w.Code, resp2.Code)
	assert.Empty(t, resp2.Message)
	assert.Equal(t, resp.DatasetId, resp2.Dataset.Id)
	assert.Equal(t, resp2.Dataset.DisplayName, "untitled-1")

	// Try to get a non-existing dataset
	req, err = http.NewRequest("GET", fmt.Sprintf("/rest/dataset/%d", rand.Uint64()), nil)
	if err != nil {
		t.Fatalf("failed to build http request with err: %v", err)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp3 manager.GetDatasetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp3); err != nil {
		t.Fatalf("failed to unmarshal json with err: %v", err)
	}

	assert.Equal(t, w.Code, http.StatusNotFound, "response code")
	assert.Equal(t, w.Code, resp3.Code)
	assert.NotEmpty(t, resp3.Message)
	assert.Empty(t, resp3.Dataset)
}

func TestListDatasets(t *testing.T) {
	router := GetRouter()
	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("POST", "/rest/dataset", nil)
		if err != nil {
			t.Fatalf("failed to build new http request with err: %v", err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusCreated)
	}

	req, err := http.NewRequest("GET", "/rest/datasets", nil)
	if err != nil {
		t.Fatalf("failed to bring new http request with err: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var resp manager.ListDatasetsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal json with err: %v", err)
	}

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Code, resp.Code)
	assert.Equal(t, resp.TotalDatasets, 10)
}

func TestUploadDataset(t *testing.T) {
	router := GetRouter()

	// Create a Dataset
	req, err := http.NewRequest("POST", "/rest/dataset", nil)
	if err != nil {
		t.Fatalf("failed to build new http request with err: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusCreated)

	var resp manager.CreateDatasetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal json with err: %v", err)
	}

	// TODO(#14): Create tests for UploadDataset
}
