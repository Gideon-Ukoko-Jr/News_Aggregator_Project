package handler

import (
	"content-delivery-service/internal/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ContentHandler struct {
}

func NewContentHandler() *ContentHandler {
	return &ContentHandler{}
}

func (ch *ContentHandler) GetUserNewsContent(c *gin.Context) {
	// Check and fetch Authorization header
	authorizationToken := c.GetHeader("Authorization")
	if authorizationToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	// Parse pagination parameters
	page, pageSize, err := ParsePaginationParameters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse categories parameter
	categoriesString := c.Query("categories")
	var categories []string
	if categoriesString != "" {
		categories = strings.Split(categoriesString, ",")
		fmt.Println("Categories:")
		for _, category := range categories {
			fmt.Println(category)
		}
	} else {
		// Fetch user preferences if no categories specified
		userPrefs, err := utils.GetUserPreferences(authorizationToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = userPrefs.Preferences
	}

	// Parse and validate keyword parameter
	keyword := c.Query("keyword")

	// Call GetRecentNews method
	newsApiResponse, err := utils.GetRecentNews(page, pageSize, keyword, categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serve the response to the user
	c.JSON(http.StatusOK, gin.H{"data": newsApiResponse})
}

func ParsePaginationParameters(c *gin.Context) (int, int, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		return 0, 0, errors.New("invalid page parameter")
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		return 0, 0, errors.New("invalid pageSize parameter")
	}

	return page, pageSize, nil
}
