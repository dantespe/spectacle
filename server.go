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

	// UI
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := handler.AddUIHandlerRoutes(router, wd); err != nil {
		log.Fatal(err)
	}

	router.Run() // localhost:8080

}
