package utils

import (
	"content-delivery-service/config"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UserPreferences struct {
	Preferences []string `json:"preferences"`
}

type NewsApiResponse struct {
	IsFirstPage    bool          `json:"IsFirstPage"`
	IsLastPage     bool          `json:"IsLastPage"`
	NewsContent    []NewsContent `json:"newsContent"`
	Page           int           `json:"page"`
	PageSize       int           `json:"pageSize"`
	Source         string        `json:"source"`
	SpecialKeyUsed bool          `json:"specialKeyUsed"`
	TotalPages     int           `json:"totalPages"`
}
type NewsContent struct {
	ID          int       `json:"ID"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	DeletedAt   time.Time `json:"DeletedAt"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Category    string    `json:"category"`
	ApiSource   string    `json:"apiSource"`
}

func GetRecentNews(page int, pageSize int, keyword string, categories []string) (*NewsApiResponse, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	newsAggregatorRecentNewsURL := cfg.NewsAggregatorRecentNewsURL

	q := url.Values{}

	if page > 0 {
		q.Add("page", strconv.Itoa(page))
	}

	if pageSize > 0 {
		q.Add("pageSize", strconv.Itoa(pageSize))
	}

	if keyword != "" {
		q.Add("keyword", keyword)
	}

	if categories != nil && len(categories) > 0 {
		q.Add("categories", strings.Join(categories, ","))
	}

	// Append query parameters to the URL
	if len(q) > 0 {
		newsAggregatorRecentNewsURL += "?" + q.Encode()
	}

	println(newsAggregatorRecentNewsURL)
	req, err := http.NewRequest("GET", newsAggregatorRecentNewsURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Special-Key", cfg.SpecialKey)

	// Make the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to fetch news content")
	}

	// Decode the response body into UserPreferences struct
	var newsApiResponse NewsApiResponse
	err = json.NewDecoder(resp.Body).Decode(&newsApiResponse)
	if err != nil {
		return nil, err
	}

	return &newsApiResponse, nil
}
func GetUserPreferences(authorizationToken string) (*UserPreferences, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	userServicePreferenceURL := cfg.UserServicePreferenceURL

	// Create a request to the user-service
	req, err := http.NewRequest("GET", userServicePreferenceURL, nil)
	if err != nil {
		return nil, err
	}

	// Set Authorization header
	req.Header.Set("Authorization", authorizationToken)

	// Make the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to fetch user preferences")
	}

	// Decode the response body into UserPreferences struct
	var userPrefs UserPreferences
	err = json.NewDecoder(resp.Body).Decode(&userPrefs)
	if err != nil {
		return nil, err
	}

	return &userPrefs, nil
}
