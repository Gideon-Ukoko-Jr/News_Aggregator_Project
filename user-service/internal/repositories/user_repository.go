package repositories

import (
	"gorm.io/gorm"
	"strings"
	"user-service/internal/models"
	"user-service/internal/utils"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Creating a new user in the db
func (ur *UserRepository) CreateUser(user *models.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return utils.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

// Retrieving a user by ID from the db
func (ur *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := ur.db.First(&user, userID).Error
	return &user, err
}

// Retrieving a user by username from the db
func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
