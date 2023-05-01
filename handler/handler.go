// Package handler contains logic from translating web requests the spectacle backend.package handler
package handler

import (
    "log"
    "net/http"
    "github.com/dantespe/spectacle/manager"
    "github.com/gin-gonic/gin"  	
)

type RestHandler struct {
    mgr *manager.Manager
    rb *manager.RequestBuilder
}

func NewRestHandler() (*RestHandler, error) {
    mgr, err := manager.New()
    if err != nil {
        return nil, err
    }

    return &RestHandler{
        mgr: mgr,
        rb: &manager.RequestBuilder{},
    }, nil
}

func (h *RestHandler) Status(c *gin.Context) {
    resp := h.mgr.Status()
    c.JSON(resp.ResponseCode(), resp.JSON())
}
  
func (h *RestHandler) CreateDataset(c *gin.Context) { 
    req, err := h.rb.CreateDatasetRequestBuilder(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{
            "message": err.Error(),
        })
    }
    resp := h.mgr.CreateDataset(req)
    c.JSON(resp.ResponseCode(), resp.JSON())
}
  
func (h *RestHandler) GetDataset(c *gin.Context) {
    req, err := h.rb.GetDatasetRequestBuilder(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{
        "message": err.Error(),
        })
        return
    }
    resp := h.mgr.GetDataset(req)
    c.JSON(resp.ResponseCode(), resp.JSON())
}
  
func (h *RestHandler) ListDatasets(c *gin.Context) {
    req, err := h.rb.ListDatasetsRequestBuilder(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{
        "message": err.Error(),
        })
        return
    }
    resp := h.mgr.ListDatasets(req)
    c.JSON(resp.ResponseCode(), resp.JSON())
}
  
func (h *RestHandler) UploadDataset(c *gin.Context)  {
    resp, err := h.rb.UploadDatasetRequestBuilder(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{
        "message": err.Error(),
        })
    }
  
    _, header, err := c.Request.FormFile("file")
    if err != nil {
      c.JSON(http.StatusBadRequest, map[string]string{
        "message": err.Error(),
      })
    }
    log.Printf("Reading filename: %s", header.Filename)

    c.JSON(200, map[string]interface{}{
      "filename": header.Filename,
      "resp": resp,
    })
}

func (h *RestHandler) GetRoutes() map[string]gin.HandlerFunc {
    return map[string]gin.HandlerFunc {
        "/status": h.Status,
        "/datasets": h.ListDatasets,
        "/dataset/:id": h.GetDataset,
    }
}

func (h *RestHandler) PostRoutes() map[string]gin.HandlerFunc {
    return map[string]gin.HandlerFunc {
        "/dataset": h.CreateDataset,
        // "/datasets/:id/upload": h.UploadDataset,
    }
}