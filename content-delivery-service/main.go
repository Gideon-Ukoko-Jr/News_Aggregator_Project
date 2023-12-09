package main

import (
	"content-delivery-service/internal/handler"
	"github.com/gin-gonic/gin"
)

const apiPrefix = "/api/v1/content-delivery-service"

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// Initialize ContentHandler
	contentHandler := handler.NewContentHandler()

	// Define routes
	router.GET(apiPrefix+"/user-news", contentHandler.GetUserNewsContent)

	// Start the server
	err := router.Run(":8082")
	if err != nil {
		panic("Failed to start the server: " + err.Error())
	}
}
