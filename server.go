// Package spectacle is the main entry point for spectacle.
package main

import (
	"log"
	"os"

	"github.com/dantespe/spectacle/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// TODO: Add UI
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var assets_folder string = wd + "/assets"
	var page_files string = wd + "/pages"
	var style_file string = wd + "/templates/layout/style.css"

	router.Static("assets", assets_folder)
	router.Static("templates/pages", page_files)
	router.StaticFile("/style.css", style_file)

	rh, err := handler.NewRestHandler()
	if err != nil {
		log.Fatal(err)
	}
	rest := router.Group("rest")
	{
		for k, v := range rh.GetRoutes() {
			rest.GET(k, v)
		}
		for k, v := range rh.PostRoutes() {
			rest.POST(k, v)
		}
	}

	uh, uiErr := handler.NewUIHandler()
	if uiErr != nil {
		log.Fatal(uiErr)
	}
	{
		for k, v := range uh.GetRoutes() {
			router.GET(k, v)
		}
		for k, v := range uh.PostRoutes() {
			router.POST(k, v)
		}
	}
	// TODO: Add UI Handler
	router.Run() // localhost:8080
}
