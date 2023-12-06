package models

import (
	"gorm.io/gorm"
	"time"
)

type NewsContent struct {
	gorm.Model
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Category    string    `json:"category"`
	ApiSource   string    `json:"apiSource"`
}
