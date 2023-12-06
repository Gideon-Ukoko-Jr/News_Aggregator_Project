package repositories

import (
	"gorm.io/gorm"
	"news-aggregator-service/internal/models"
)

type NewsContentRepository struct {
	db *gorm.DB
}

func NewNewsContentRepository(db *gorm.DB) *NewsContentRepository {
	return &NewsContentRepository{db: db}
}

// Saving news content to the database after validating
func (nr *NewsContentRepository) SaveNewsContent(newsContent *models.NewsContent) error {
	if err := nr.db.Create(newsContent).Error; err != nil {
		return err
	}
	return nil
}

// Checking if news content with a similar title and published date exists
func (nr *NewsContentRepository) NewsContentExists(newsContent *models.NewsContent) bool {
	var count int64

	// Check for similar titles and published date within 24 hours
	nr.db.Model(&models.NewsContent{}).
		Where("levenshtein(title, ?) <= ? AND ABS(EXTRACT(EPOCH FROM (published_at - ?))) <= 86400",
			newsContent.Title, 4, newsContent.PublishedAt).
		Count(&count)

	return count > 0
}

func (ncr *NewsContentRepository) GetPaginatedNewsContent(page, pageSize int) ([]models.NewsContent, int64, error) {
	var total int64
	var newsContent []models.NewsContent

	// Get total count
	if err := ncr.db.Model(&models.NewsContent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query paginated news content with sorting
	if err := ncr.db.Offset(offset).Limit(pageSize).Order("updated_at desc").Find(&newsContent).Error; err != nil {
		return nil, 0, err
	}

	return newsContent, total, nil
}
