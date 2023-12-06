// services/news_service.go

package services

import (
	"fmt"
	"news-aggregator-service/internal/models"
	"news-aggregator-service/internal/repositories"
	"news-aggregator-service/internal/utils"
	"strings"
	"time"
)

type NewsService struct {
	newsContentRepository *repositories.NewsContentRepository
}

func NewNewsService(newsContentRepository *repositories.NewsContentRepository) *NewsService {
	return &NewsService{newsContentRepository: newsContentRepository}
}

// UpdateNewsContentFromNewsAPI fetches news content from the News API and saves it to the database.
func (ns *NewsService) UpdateNewsContentFromNewsAPI() {
	config := utils.NewsAPI

	for _, category := range config.Categories {
		apiURL := strings.Replace(config.URL, "{apiKey}", config.Key, 1)
		apiURL = strings.Replace(apiURL, "{category}", category, 1)

		response := ns.fetchNewsContent(apiURL)

		if response != nil && response.Status == "ok" {
			for _, article := range response.Articles {
				ns.saveNewsContent(article, category)
			}
		}
	}
}

// UpdateNewsContentFromGuardianNewsAPI fetches news content from the Guardian News API and saves it to the database.
func (ns *NewsService) UpdateNewsContentFromGuardianNewsAPI() {
	config := utils.GuardianNewsAPI

	for _, section := range config.Sections {
		// Generate today's date in yyyy-mm-dd format
		today := time.Now().Format("2006-01-02")

		// Replace placeholders in the URL with actual values
		apiURL := strings.Replace(config.URL, "{apiKey}", config.Key, 1)
		apiURL = strings.Replace(apiURL, "{sectionName}", section, 1)
		apiURL = strings.Replace(apiURL, "{fromDate}", today, 1)
		apiURL = strings.Replace(apiURL, "{toDate}", today, 1)

		response := ns.fetchGuardianNewsContent(apiURL)

		if response != nil && response.Response.Status == "ok" {
			for _, result := range response.Response.Results {
				ns.saveGuardianNewsContent(result, section)
			}
		}
	}
}

func (ns *NewsService) fetchNewsContent(apiURL string) *models.NewsApiResponse {
	// Implement the logic to fetch news content from the News API
	response, err := utils.GetNewsAPIResponse(apiURL)

	// Log errors if any
	if err != nil {
		fmt.Printf("Error fetching news content from News API: %s\n", apiURL)
		return nil
	}

	return response
}

// Fetching news content from the Guardian News API.
func (ns *NewsService) fetchGuardianNewsContent(apiURL string) *models.GuardianNewsApiResponse {
	// Implement the logic to fetch news content from the Guardian News API
	response, err := utils.GetGuardianNewsAPIResponse(apiURL)

	// Log errors if any
	if err != nil {
		fmt.Printf("Error fetching news content from Guardian News API: %s\n", apiURL)
		return nil
	}

	return response
}

// saveNewsContent validates and saves news content from the News API to the database.
func (ns *NewsService) saveNewsContent(article models.ArticleApiResponse, category string) {
	// Implementing the logic to validate and save news content to the database
	newsContent := models.NewsContent{
		Author:      article.Author,
		Title:       article.Title,
		URL:         article.URL,
		URLToImage:  article.URLToImage,
		PublishedAt: article.PublishedAt,
		Category:    category,
		ApiSource:   "NEWS API",
	}

	// Saving to the database
	err := ns.newsContentRepository.SaveNewsContent(&newsContent)
	if err != nil {
		fmt.Printf("Error saving news content from News API to the database: %v\n", err)
	}
}

// Validating and saveing news content from the Guardian News API to the database.
func (ns *NewsService) saveGuardianNewsContent(result models.GuardianResult, section string) {
	// Implement the logic to validate and save news content to the database
	newsContent := models.NewsContent{
		Author:      result.WebTitle,
		Title:       result.WebTitle,
		URL:         result.WebUrl,
		URLToImage:  result.Fields.Thumbnail,
		PublishedAt: utils.ParseGuardianDate(result.WebPublicationDate),
		Category:    section,
		ApiSource:   "GUARDIAN API",
	}

	// Saving to the database
	err := ns.newsContentRepository.SaveNewsContent(&newsContent)
	if err != nil {
		fmt.Printf("Error saving news content from Guardian News API to the database: %v\n", err)
	}
}
