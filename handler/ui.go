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

func (u *UIHandler) CreateChart(c *gin.Context) {
	c.HTML(http.StatusOK, "build_charts.html", nil)
}

func (u *UIHandler) CreateDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "build_dashboards.html", nil)
}

func (u *UIHandler) CreateSummary(c *gin.Context) {
	c.HTML(http.StatusOK, "build_summary.html", nil)
}

func (u *UIHandler) ImportData(c *gin.Context) {
	c.HTML(http.StatusOK, "view_datasets.html", nil)
}

func (u *UIHandler) Share(c *gin.Context) {
	c.HTML(http.StatusOK, "sharing.html", nil)
}

func (u *UIHandler) Settings(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.html", nil)
}

func (u *UIHandler) Visualizations(c *gin.Context) {
	c.HTML(http.StatusOK, "visualizations.html", nil)
}

func (u *UIHandler) GetRoutes() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"/":                 u.Index,
		"/get_started":      u.Starter,
		"/create_chart":     u.CreateChart,
		"/create_dashboard": u.CreateDashboard,
		"/create_summary":   u.CreateSummary,
		"/create_dataset":   u.ImportData,
		"/visualizations":   u.Visualizations,
		"/share":            u.Share,
		"/settings":         u.Settings,
	}
}

func (u *UIHandler) PostRoutes() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		// "/create_chart": ...
	}
}
