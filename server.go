// Package spectacle is the main entry point for spectacle.
package main

import (
  "log"

  "github.com/dantespe/spectacle/handler"
  "github.com/gin-gonic/gin"
)

func main() {
    // TODO: Add UI 
    router := gin.Default()
    
    rh, err := handler.NewRestHandler()
    if err != nil {
        log.Fatal(err)
    }
    rest := router.Group("rest")
    {
        for k,v := range rh.GetRoutes() {
            rest.GET(k, v)
        } 
        for k,v := range rh.PostRoutes() {
            rest.POST(k, v)
        }
    }
    // TODO: Add UI Handler 
    router.Run() // localhost:8080
}