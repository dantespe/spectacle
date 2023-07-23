package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
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

func (u *UIHandler) MakeDataset(c *gin.Context) {
	datasetName, _ := c.MultipartForm()
	var fileName string

	if len(datasetName.Value["newFileName"]) == 0 {
		fileName = fmt.Sprintf(`{"displayName": "%s", "hasHeaders": true}`, "NewDataset")
	} else {
		fileName = fmt.Sprintf(`{"displayName": "%s", "hasHeaders": true}`, datasetName.Value["newFileName"][0])
	}

	datasetResp, err := http.Post(
		"http://localhost:8080/rest/dataset",
		"application/json", strings.NewReader(fileName))
	if err != nil {
		log.Printf("Error with creating new dataset: %v", err)
	}
	defer datasetResp.Body.Close()

	if datasetResp.StatusCode == http.StatusCreated {
		var resp manager.CreateDatasetResponse
		bodyBytes, err := io.ReadAll(datasetResp.Body)
		if err != nil {
			log.Print("Processing read error: ", err)
		}
		json.Unmarshal(bodyBytes, &resp)

		file, _, err := c.Request.FormFile("file")
		if err != nil {
			log.Print("Processing form error: ", err)
		}

		client := http.Client{}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fw, err := writer.CreateFormFile("file", "dataset.csv")
		if err != nil {
			log.Printf("There was an issue creating a new field: %v", err)
		}
		_, err = io.Copy(fw, file)
		if err != nil {
			log.Printf("Failed to build HTTP request with err: %v", err)
		}

		writer.Close()
		uploadResp, err := http.NewRequest("POST", "http://localhost:8080/rest"+resp.DatasetUrl+"/upload", bytes.NewReader(body.Bytes()))
		if err != nil {
			log.Printf("There's a bug in recreating a new request: %v", err)
		}
		uploadResp.Header.Set("Content-Type", writer.FormDataContentType())

		client.Do(uploadResp)
		c.Redirect(http.StatusFound, "http://localhost:8080/edit_dataset")
	}
}

func (u *UIHandler) CreateChart(c *gin.Context) {
	c.HTML(http.StatusOK, "chart_builder.html", nil)
}

func (u *UIHandler) CreateTemplate(c *gin.Context) {
	c.HTML(http.StatusOK, "template_builder.html", nil)
}

func (u *UIHandler) CreateDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard_builder.html", nil)
}

func (u *UIHandler) CreateSummary(c *gin.Context) {
	c.HTML(http.StatusOK, "summary_builder.html", nil)
}

func (u *UIHandler) ImportData(c *gin.Context) {
	c.HTML(http.StatusCreated, "import.html", nil)
}

func (u *UIHandler) EditData(c *gin.Context) {
	datasetsResp, err := http.Get("http://localhost:8080/rest/datasets")
	if err != nil {
		log.Printf("Error occured while trying to retrieve datasets because %s", err)
	}

	var datasets manager.ListDatasetsResponse
	newBodyBytes, err := io.ReadAll(datasetsResp.Body)
	if err != nil {
		log.Print("Processing read error: ", err)
	}
	json.Unmarshal(newBodyBytes, &datasets)

	c.HTML(http.StatusOK, "edit_dataset.html", gin.H{
		"datasets": datasets.Results,
	})
}

func (u *UIHandler) Settings(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.html", nil)
}

func (u *UIHandler) Visualize(c *gin.Context) {
	c.HTML(http.StatusOK, "visualize.html", nil)
}

func (u *UIHandler) Studio(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "studio.html", gin.H{
		"title":      "Data Studio",
		"img_source": "/assets/images/spectacle.png",
		"datasetId":  id,
	})
}

func (u *UIHandler) GetRoutes() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"/":                 u.Index,
		"/create_chart":     u.CreateChart,
		"/create_template":  u.CreateTemplate,
		"/create_dashboard": u.CreateDashboard,
		"/create_summary":   u.CreateSummary,
		"/create_dataset":   u.ImportData,
		"/edit_dataset":     u.EditData,
		"/visualize":        u.Visualize,
		"/settings":         u.Settings,
		"/studio/:id":       u.Studio,
	}
}

func (u *UIHandler) PostRoutes() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"/create_dataset": u.MakeDataset,
		// "/edit_dataset":   u.EditData,
	}
}
