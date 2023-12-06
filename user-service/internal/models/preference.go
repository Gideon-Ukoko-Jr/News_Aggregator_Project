package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Preference struct {
	gorm.Model
	Username   string         `gorm:"unique_index;not null;unique" json:"username" binding:"required"`
	Categories pq.StringArray `json:"categories" gorm:"type:text[]" binding:"required"`
}

type PreferencesRequest struct {
	Categories pq.StringArray `json:"categories" binding:"required"`
}
