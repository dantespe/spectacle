// Package spectacle is the main entry point for spectacle.
package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func StatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "HEALTHY",
	})
}

func CreateDatasetHandler(c *gin.Context) {
}


func main() {
  r := gin.Default()

  r.GET("/status", StatusHandler)

  r.Run() // localhost:8080
}