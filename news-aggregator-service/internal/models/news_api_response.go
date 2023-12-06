package models

import "time"

type NewsApiResponse struct {
	Status       string               `json:"status"`
	TotalResults int                  `json:"totalResults"`
	Articles     []ArticleApiResponse `json:"articles"`
}

type ArticleApiResponse struct {
	Source      SourceApiResponse `json:"source"`
	Author      string            `json:"author"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	URL         string            `json:"url"`
	URLToImage  string            `json:"urlToImage"`
	PublishedAt time.Time         `json:"publishedAt"`
}

type SourceApiResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
