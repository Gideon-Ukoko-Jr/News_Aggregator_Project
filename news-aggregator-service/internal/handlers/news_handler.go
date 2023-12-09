// handlers/news_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"news-aggregator-service/internal/models"
	_ "news-aggregator-service/internal/models"
	"news-aggregator-service/internal/repositories"
	"news-aggregator-service/internal/utils"
	"strconv"
	"strings"
	"time"
)

const (
	internalServerError  = "Internal Server Error"
	invalidPageError     = "Invalid page or pageSize"
	invalidSpecialKey    = "Invalid Special Key"
	twelveHoursInSeconds = 12 * 60 * 60
)

type NewsHandler struct {
	newsContentRepository *repositories.NewsContentRepository
	redisClient           *redis.Client
}

// NewNewsHandler creates a new NewsHandler instance.
func NewNewsHandler(newsContentRepository *repositories.NewsContentRepository, redisClient *redis.Client) *NewsHandler {
	return &NewsHandler{
		newsContentRepository: newsContentRepository,
		redisClient:           redisClient,
	}
}

func (nh *NewsHandler) GetNewsContentHandler(c *gin.Context) {
	// Parsing query parameters
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))

	// Validations
	if page <= 0 || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidPageError})
		return
	}

	// Fetching paginated news content from the repository
	newsContent, total, err := nh.newsContentRepository.GetPaginatedNewsContent(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalServerError})
		return
	}

	// Calculating total number of pages
	totalPages := utils.CalculateTotalPages(total, int64(pageSize))

	// Determining if it's the first and last page
	isFirstPage := page == 1
	isLastPage := page == totalPages

	// Preparing the response
	response := gin.H{
		"page":        page,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
		"isFirstPage": isFirstPage,
		"isLastPage":  isLastPage,
		"content":     newsContent,
	}

	// Returning the response
	c.JSON(http.StatusOK, response)
}

func (nh *NewsHandler) GetPaginatedNewsContentFilteredHandler(c *gin.Context) {
	// Validating special key header
	if !isValidSpecialKey(c.Request.Header.Get("Special-Key")) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidSpecialKey})
		return
	}

	// Parsing and validating pagination parameters
	page, pageSize, err := utils.ParsePaginationParameters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parsing and validate categories parameter
	categoriesString := c.Query("categories")
	var categories []string
	if categoriesString != "" {
		categories = strings.Split(categoriesString, ",")
		fmt.Println("Categories:")
		for _, category := range categories {
			fmt.Println(category)
		}
	} else {
		categories = []string{} // Ensure an empty slice if no categories are provided
		fmt.Println("Categories: (none)")
	}

	// Parse and validating keyword parameter
	keyword := c.Query("keyword")

	// Parsing and validating publishedAfter parameter
	publishedAfter, err := utils.ParseTimeParameter(c, "publishedAfter")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetching paginated and filtered news content from the repository
	newsContent, total, err := nh.newsContentRepository.GetPaginatedNewsContentFiltered(page, pageSize, categories, keyword, publishedAfter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news content"})
		return
	}

	// Responding with the paginated and filtered news content
	c.JSON(http.StatusOK, gin.H{
		"page":           page,
		"pageSize":       pageSize,
		"total":          total,
		"newsContent":    newsContent,
		"specialKeyUsed": true,
	})
}

//func (nh *NewsHandler) GetRecentNewsHandler(c *gin.Context) {
//
//	if !isValidSpecialKey(c.Request.Header.Get("Special-Key")) {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidSpecialKey})
//		return
//	}
//
//	// Calculate the time six hours ago from the current time
//	publishedAfter := time.Now().Add(-time.Duration(sixHoursInSeconds) * time.Second)
//
//	println(publishedAfter.GoString())
//
//	// Fetch all news content published within the last six hours
//	recentNews, err := nh.newsContentRepository.GetRecentNewsContent(publishedAfter)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": internalServerError})
//		return
//	}
//
//	// Respond with the recent news content
//	c.JSON(http.StatusOK, gin.H{
//		"recentNews": recentNews,
//	})
//}

//func (nh *NewsHandler) GetRecentNewsHandler(c *gin.Context) {
//	if !isValidSpecialKey(c.Request.Header.Get("Special-Key")) {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidSpecialKey})
//		return
//	}
//
//	// Check if recent news is in the cache
//	recentNews, err := nh.getRecentNewsFromCache()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": internalServerError})
//		return
//	}
//
//	// If cache is empty, fetch news from the database
//	if recentNews == nil {
//		publishedAfter := time.Now().Add(-time.Duration(sixHoursInSeconds) * time.Second)
//		recentNews, err = nh.newsContentRepository.GetRecentNewsContent(publishedAfter)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": internalServerError})
//			return
//		}
//
//		// Update the cache with the fetched news
//		err := nh.updateRecentNewsCache(recentNews)
//		if err != nil {
//			fmt.Printf("Error updating recent news cache: %v\n", err)
//		}
//	}
//
//	// Respond with the recent news
//	c.JSON(http.StatusOK, gin.H{"recentNews": recentNews})
//}

//func (nh *NewsHandler) GetRecentNewsHandler(c *gin.Context) {
//	if !isValidSpecialKey(c.Request.Header.Get("Special-Key")) {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidSpecialKey})
//		return
//	}
//
//	// Parse and validate query parameters
//	page, pageSize, err := utils.ParsePaginationParameters(c)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	// Parse and validate categories parameter
//	categoriesString := c.Query("categories")
//	var categories []string
//	if categoriesString != "" {
//		categories = strings.Split(categoriesString, ",")
//		fmt.Println("Categories:")
//		for _, category := range categories {
//			fmt.Println(category)
//		}
//	} else {
//		categories = []string{} // Ensure an empty slice if no categories are provided
//		fmt.Println("Categories: (none)")
//	}
//
//	// Parse and validate keyword parameter
//	keyword := c.Query("keyword")
//
//	// Check if recent news is in the cache
//	recentNews, err := nh.getRecentNewsFromCache()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": internalServerError})
//		return
//	}
//
//	// If cache is not empty, filter based on categories and keyword
//	if recentNews != nil {
//		recentNews = filterRecentNews(recentNews, categories, keyword)
//	}
//
//	// If cache is still empty or after filtering, fetch news from the database
//	if recentNews == nil {
//		publishedAfter := time.Now().Add(-time.Duration(sixHoursInSeconds) * time.Second)
//		recentNews, err = nh.newsContentRepository.GetRecentNewsContent(publishedAfter)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": internalServerError})
//			return
//		}
//
//		// Update the cache with the fetched news
//		err := nh.updateRecentNewsCache(recentNews)
//		if err != nil {
//			fmt.Printf("Error updating recent news cache: %v\n", err)
//		}
//	}
//
//	// Respond with the paginated and filtered recent news
//	c.JSON(http.StatusOK, gin.H{
//		"page":           page,
//		"pageSize":       pageSize,
//		"newsContent":    recentNews,
//		"specialKeyUsed": true,
//	})
//}

func (nh *NewsHandler) GetRecentNewsHandler(c *gin.Context) {
	if !isValidSpecialKey(c.Request.Header.Get("Special-Key")) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidSpecialKey})
		return
	}

	// Parse and validate query parameters
	page, pageSize, err := utils.ParsePaginationParameters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse and validate categories parameter
	categoriesString := c.Query("categories")
	var categories []string
	if categoriesString != "" {
		categories = strings.Split(categoriesString, ",")
		fmt.Println("Categories:")
		for _, category := range categories {
			fmt.Println(category)
		}
	} else {
		categories = []string{} // Ensure an empty slice if no categories are provided
		fmt.Println("Categories: (none)")
	}

	// Parse and validate keyword parameter
	keyword := c.Query("keyword")

	// Checking if pagination is requested
	isPaginated := page > 0 && pageSize > 0

	// Checking if recent news is in the cache
	recentNews, err := nh.getRecentNewsFromCache()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalServerError})
		return
	}

	// Filtering based on categories and keyword if cache is not empty
	if recentNews != nil {
		recentNews = filterRecentNews(recentNews, categories, keyword)
	}

	// Fetching news from the database if cache is empty
	if recentNews == nil {
		publishedAfter := time.Now().Add(-time.Duration(twelveHoursInSeconds) * time.Second)
		recentNews, err = nh.newsContentRepository.GetRecentNewsContent(publishedAfter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": internalServerError})
			return
		}

		// Updating the cache with the fetched news
		err := nh.updateRecentNewsCache(recentNews)
		if err != nil {
			fmt.Printf("Error updating recent news cache: %v\n", err)
		}
	}

	// Responding with paginated or non-paginated recent news based on the request
	if isPaginated {
		// Paginated response
		totalPages := utils.CalculateTotalPages(int64(len(recentNews)), int64(pageSize))
		isFirstPage := page == 1
		isLastPage := page == totalPages

		c.JSON(http.StatusOK, gin.H{
			"page":           page,
			"pageSize":       pageSize,
			"totalPages":     totalPages,
			"isFirstPage":    isFirstPage,
			"isLastPage":     isLastPage,
			"newsContent":    getPaginatedNews(recentNews, page, pageSize),
			"source":         getSourceInfo(recentNews),
			"specialKeyUsed": true,
		})
	} else {
		// Non-paginated response
		c.JSON(http.StatusOK, gin.H{
			"newsContent": getNonPaginatedNews(recentNews),
			"source":      getSourceInfo(recentNews),
		})
	}
}

func getSourceInfo(news []models.NewsContent) string {
	if len(news) > 0 {
		return "Cache"
	}
	return "Database"
}

func getPaginatedNews(news []models.NewsContent, page, pageSize int) []models.NewsContent {
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= len(news) {
		return []models.NewsContent{}
	}

	if endIndex > len(news) {
		endIndex = len(news)
	}

	return news[startIndex:endIndex]
}

// getNonPaginatedNews returns the entire set of news without pagination
func getNonPaginatedNews(news []models.NewsContent) []models.NewsContent {
	return news
}

// getRecentNewsFromCache retrieves recent news from the Redis cache.
func (nh *NewsHandler) getRecentNewsFromCache() ([]models.NewsContent, error) {
	ctx := context.Background()
	cacheKey := "recentNews"

	// Check if recent news is in the cache
	cachedData, err := nh.redisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Key does not exist in cache
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Unmarshal cached data
	var recentNews []models.NewsContent
	err = json.Unmarshal([]byte(cachedData), &recentNews)
	if err != nil {
		return nil, err
	}

	return recentNews, nil
}

// updateRecentNewsCache updates the Redis cache with recent news.
func (nh *NewsHandler) updateRecentNewsCache(recentNews []models.NewsContent) error {
	ctx := context.Background()
	cacheKey := "recentNews"

	// Marshal news data
	cachedData, err := json.Marshal(recentNews)
	if err != nil {
		return err
	}

	// Set the cache with the recent news
	err = nh.redisClient.Set(ctx, cacheKey, cachedData, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func filterRecentNews(news []models.NewsContent, categories []string, keyword string) []models.NewsContent {
	var filteredNews []models.NewsContent

	for _, n := range news {
		// Check if news matches categories or keyword
		if ((categories != nil && len(categories) == 0) || containsCategory(categories, n.Category)) &&
			(keyword == "" || strings.Contains(strings.ToLower(n.Title), strings.ToLower(keyword))) {
			filteredNews = append(filteredNews, n)
		}
	}

	return filteredNews
}

// containsCategory checks if a category is present in the provided list
func containsCategory(categories []string, category string) bool {
	for _, c := range categories {
		if c == category {
			return true
		}
	}
	return false
}

func isValidSpecialKey(specialKey string) bool {
	return specialKey == utils.SpecialKey
}
