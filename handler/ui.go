package handler

import (
    "net/http"

    "github.com/dantespe/spectacle/manager"
    "github.com/gin-gonic/gin"
)

type UIHandler struct {
    mgr *manager.Manager
    rb  *manager.RequestBuilder
}

func NewUIHandler() (*UIHandler, error) {
    mgr, err := manager.New()
    if err != nil {
        return nil, err
    }

    return &UIHandler{
        mgr: mgr,
        rb:  &manager.RequestBuilder{},
    }, nil
}

func (u *UIHandler) Index(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{
        "img_source": "/assets/images/spectacle.png",
    })
}

func (u *UIHandler) Starter(c *gin.Context) {
    c.HTML(http.StatusOK, "get_started.html", nil)
}

// func Home() any {
// 	return map[string]int{"datasetId": 9908576756734740233}
// }

func (u *UIHandler) GetRoutes() map[string]gin.HandlerFunc {
    return map[string]gin.HandlerFunc{
        "/welcome":     u.Index,
        "/get_started": u.Starter,
        // "/login":    u.ListDatasets,
        // "/settings": u.GetDataset,
    }
}

func (u *UIHandler) PostRoutes() map[string]gin.HandlerFunc {
    return map[string]gin.HandlerFunc{
        // "/create_chart": ...
    }
}
