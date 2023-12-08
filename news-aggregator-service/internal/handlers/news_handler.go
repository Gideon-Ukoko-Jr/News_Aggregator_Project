// handlers/news_handler.go
package handlers

import (
	"fmt"
	"net/http"
	_ "news-aggregator-service/internal/models"
	"news-aggregator-service/internal/repositories"
	"news-aggregator-service/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	internalServerError = "Internal Server Error"
	invalidPageError    = "Invalid page or pageSize"
	invalidSpecialKey   = "Invalid Special Key"
)

type NewsHandler struct {
	newsContentRepository *repositories.NewsContentRepository
	specialKey            string
}

func NewNewsHandler(newsContentRepository *repositories.NewsContentRepository) *NewsHandler {
	return &NewsHandler{
		newsContentRepository: newsContentRepository,
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
	categories := strings.Split(categoriesString, ",")
	fmt.Println("Categories:")
	for _, category := range categories {
		fmt.Println(category)
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

func isValidSpecialKey(specialKey string) bool {
	return specialKey == utils.SpecialKey
}
