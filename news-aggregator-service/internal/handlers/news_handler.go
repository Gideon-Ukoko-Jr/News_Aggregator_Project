// handlers/news_handler.go
package handlers

import (
	"news-aggregator-service/internal/repositories"
	"news-aggregator-service/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	newsContentRepository *repositories.NewsContentRepository
}

func NewNewsHandler(newsContentRepository *repositories.NewsContentRepository) *NewsHandler {
	return &NewsHandler{newsContentRepository: newsContentRepository}
}

// request to fetch all news content with pagination
func (nh *NewsHandler) GetNewsContentHandler(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))

	// Set default values if not provided
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// Fetching paginated news content from the repository
	newsContent, total, err := nh.newsContentRepository.GetPaginatedNewsContent(page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
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
	c.JSON(200, response)
}
