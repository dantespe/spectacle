package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dantespe/spectacle/manager"
	"github.com/gin-gonic/gin"
)

type UIHandler struct {
	mgr *manager.Manager
	rb  *manager.RequestBuilder
}

func AddUIHandlerRoutes(r *gin.Engine, wd string) error {
	r.LoadHTMLGlob("templates/**/*")

	var assets_folder string = wd + "/assets"
	var page_files string = wd + "/pages"
	var style_file string = wd + "/templates/layout/style.css"

	r.Static("assets", assets_folder)
	r.Static("templates/pages", page_files)
	r.StaticFile("/style.css", style_file)

	mgr, err := manager.New()
	if err != nil {
		return err
	}
	ui := &UIHandler{
		mgr: mgr,
		rb:  &manager.RequestBuilder{},
	}
	for k, v := range ui.GetRoutes() {
		r.GET(k, v)
	}
	for k, v := range ui.PostRoutes() {
		r.POST(k, v)
	}
	return nil
}

func (u *UIHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"img_source": "/assets/images/spectacle.png",
	})
}

func (u *UIHandler) Starter(c *gin.Context) {
	c.HTML(http.StatusOK, "get_started.html", nil)
}

func (u *UIHandler) MakeDataset(c *gin.Context) {
	// fileName := c.PostForm("datasetFile")
	datasetResp, err := http.Post(
		"http://localhost:8080/rest/dataset",
		"application/json", strings.NewReader(`{"displayName": "NewDataset"}`))
	if err != nil {
		log.Fatal(err)
	}
	defer datasetResp.Body.Close()
	log.Print(datasetResp.StatusCode == http.StatusCreated, datasetResp.StatusCode)
	if datasetResp.StatusCode == http.StatusCreated {
		bodyBytes, err := io.ReadAll(datasetResp.Body)
		if err != nil {
			log.Print("Processing read error: ", err)
		}

		var resp manager.CreateDatasetResponse
		json.Unmarshal(bodyBytes, &resp)

		file_data, err := c.FormFile("file")
		if err != nil {
			log.Print("Processing form error: ", err)
		}
		file, _ := file_data.Open()
		var filename string = "../../Downloads/" + file_data.Filename
		fileInfo, err := os.Stat(filename)
		// check if error is "file not exists"
		if err != nil {
			log.Print("Error processing", fileInfo, ": ", err)
		}
		defer file.Close()
		http.Post("/rest"+resp.DatasetUrl+"/upload", "multipart/form-data", file)
		c.Redirect(http.StatusTemporaryRedirect, "/edit_dataset")
	}
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
	c.HTML(http.StatusCreated, "datasets.html", nil)
}

func (u *UIHandler) EditData(c *gin.Context) {
	c.HTML(http.StatusOK, "datasets.html#v-pills-messages-tab", nil)
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
		"/edit_dataset":     u.EditData,
		"/visualizations":   u.Visualizations,
		"/share":            u.Share,
		"/settings":         u.Settings,
	}
}

func (u *UIHandler) PostRoutes() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"/create_dataset": u.MakeDataset,
		"/edit_dataset":   u.EditData,
	}
}
