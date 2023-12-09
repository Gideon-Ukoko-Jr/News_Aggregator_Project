// main.go

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"log"
	"news-aggregator-service/internal/handlers"
	"news-aggregator-service/internal/repositories"
	"news-aggregator-service/internal/services"
	"news-aggregator-service/internal/utils"
	_ "time"
)

const apiPrefix = "/api/v1/news-aggregator"

func main() {
	// Initialize configurations and database connection
	utils.InitConfig()
	db := utils.InitDB()

	utils.InitSpecialKey()

	// Initializing repositories and services
	newsContentRepository := repositories.NewNewsContentRepository(db)
	newsService := services.NewNewsService(newsContentRepository)

	// Running News API task on application startup
	newsService.UpdateNewsContentFromNewsAPI()

	// Scheduling News API task to run every 6 hours
	newsAPICron := cron.New()
	newsAPICron.AddFunc("@every 6h", func() {
		fmt.Println("Running News API task...")
		newsService.UpdateNewsContentFromNewsAPI()
	})
	newsAPICron.Start()

	// Running Guardian News API task on application startup
	newsService.UpdateNewsContentFromGuardianNewsAPI()

	// Scheduling Guardian News API task to run every 7 hours
	guardianAPICron := cron.New()
	guardianAPICron.AddFunc("@every 7h", func() {
		fmt.Println("Running Guardian News API task...")
		newsService.UpdateNewsContentFromGuardianNewsAPI()
	})
	guardianAPICron.Start()

	// Initializing Gin router
	router := gin.Default()

	// Initializing test handler for fetching news content
	newsHandler := handlers.NewNewsHandler(newsContentRepository)

	// Endpoint to fetch paginated news content
	router.GET(apiPrefix+"/news", newsHandler.GetNewsContentHandler)

	router.GET(apiPrefix+"/news/filtered", newsHandler.GetPaginatedNewsContentFilteredHandler)

	router.GET(apiPrefix+"/news/recent", newsHandler.GetRecentNewsHandler)

	// Running on port
	port := ":8081"
	go func() {
		if err := router.Run(port); err != nil {
			log.Fatal(err)
		}
	}()

	// Keeping the application running
	select {}
}
