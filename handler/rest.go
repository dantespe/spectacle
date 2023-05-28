// Package handler contains rest.go which implements spectacles Rest API.
package handler

import (
	"log"
	"net/http"

	"github.com/dantespe/spectacle/manager"
	"github.com/gin-gonic/gin"
)

type RestHandler struct {
	mgr *manager.Manager
	rb  *manager.RequestBuilder
}

func NewRestHandler() *RestHandler {
	mgr, _ := manager.New()
	return &RestHandler{
		mgr: mgr,
		rb:  &manager.RequestBuilder{},
	}
}

func AddRestHandlerRoutes(rg *gin.RouterGroup) error {
	rh := NewRestHandler()
	for k, v := range rh.GetRoutes() {
		rg.GET(k, v)
	}
	for k, v := range rh.PostRoutes() {
		rg.POST(k, v)
	}
	return nil
}

func (h *RestHandler) Status(c *gin.Context) {
	c.JSON(h.mgr.Status())
}

func (h *RestHandler) CreateDataset(c *gin.Context) {
	req, resp := h.rb.CreateDatasetRequestBuilder(c)
	if resp != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(h.mgr.CreateDataset(req))
}

func (h *RestHandler) GetDataset(c *gin.Context) {
	req, err := h.rb.GetDatasetRequestBuilder(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(h.mgr.GetDataset(req))
}

func (h *RestHandler) ListDatasets(c *gin.Context) {
	req, err := h.rb.ListDatasetsRequestBuilder(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(h.mgr.ListDatasets(req))
}

func (h *RestHandler) UploadDataset(c *gin.Context) {
	req, err := h.rb.UploadDatasetRequestBuilder(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(h.mgr.UploadDataset(req))
}

func (h *RestHandler) GetRoutes() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"/status":      h.Status,
		"/datasets":    h.ListDatasets,
		"/dataset/:id": h.GetDataset,
	}
}

func (h *RestHandler) PostRoutes() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"/dataset":            h.CreateDataset,
		"/dataset/:id/upload": h.UploadDataset,
	}
}
