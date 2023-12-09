package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"news-aggregator-service/internal/models"
	"time"
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

	// Check for similar titles (exact match or Levenshtein distance within 4) and published date within 24 hours
	nr.db.Model(&models.NewsContent{}).
		Where("(title = ? OR levenshtein(title, ?) <= ?) AND ABS(EXTRACT(EPOCH FROM (published_at - ?))) <= 86400",
			newsContent.Title, newsContent.Title, 4, newsContent.PublishedAt).
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

func (ncr *NewsContentRepository) GetPaginatedNewsContentFiltered(page, pageSize int, categories []string, keyword string, publishedAfter time.Time) ([]models.NewsContent, int64, error) {
	var total int64
	var newsContent []models.NewsContent

	query := ncr.db.Model(&models.NewsContent{})

	fmt.Println("Categories:", categories)
	fmt.Println("Categories: Length", len(categories))
	fmt.Println("Keyword:", keyword)
	fmt.Println("Published After:", publishedAfter)

	if categories != nil && len(categories) > 0 {
		query = query.Where("category IN (?)", categories)
	}

	if keyword != "" {
		query = query.Where("title ILIKE ?", "%"+keyword+"%")
	}

	if !publishedAfter.IsZero() {
		query = query.Where("published_at >= ?", publishedAfter)
	}
	// Getting total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculating offset
	offset := (page - 1) * pageSize

	// Querying paginated news content with sorting
	if err := query.Offset(offset).Limit(pageSize).Order("published_at desc").Find(&newsContent).Error; err != nil {
		return nil, 0, err
	}

	return newsContent, total, nil
}

func (ncr *NewsContentRepository) GetRecentNewsContent(publishedAfter time.Time) ([]models.NewsContent, error) {
	var newsContent []models.NewsContent

	query := ncr.db.Model(&models.NewsContent{})
	err := query.Where("published_at >= ?", publishedAfter).Order("published_at desc").Find(&newsContent).Error
	if err != nil {
		return nil, err
	}

	return newsContent, nil
}
