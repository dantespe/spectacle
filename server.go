// Package spectacle is the main entry point for spectacle.
package main

import (
    "log"
    "os"

    "github.com/dantespe/spectacle/handler"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // REST 
    if err := handler.AddRestHandlerRoutes(router.Group("rest")); err != nil {
        log.Fatal(err)        
    }

    // TODO: Clean this up
    router.LoadHTMLGlob("templates/*")

    wd, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    var assets_folder string = wd + "/assets"

    router.Static("assets", assets_folder)

    uh, uiErr := handler.NewUIHandler()
    if uiErr != nil {
        log.Fatal(uiErr)
    }
    ui := router.Group("home")
    {
        for k, v := range uh.GetRoutes() {
            ui.GET(k, v)
        }
        for k, v := range uh.PostRoutes() {
            ui.POST(k, v)
        }
    }

    router.Run() // localhost:8080
}
