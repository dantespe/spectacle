// Package spectacle is the main entry point for spectacle.
package main

import (
  "log"

  "github.com/dantespe/spectacle/handler"
  "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    if err := handler.AddRestHandlerRoutes(router.Group("rest")); err != nil {
        log.Fatal(err)        
    }

    router.Run() // localhost:8080
}