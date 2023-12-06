package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"user-service/internal/models"
)

type PreferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) *PreferenceRepository {
	return &PreferenceRepository{db: db}
}

func (pr *PreferenceRepository) FindByUsername(username string) (*models.Preference, error) {
	var preference models.Preference
	if err := pr.db.Where("username = ?", username).First(&preference).Error; err != nil {
		fmt.Printf("Error fetching preferences for user %s: %s\n", username, err.Error())
		return nil, err
	}
	fmt.Printf("Preferences for user %s fetched successfully: %+v\n", username, preference)
	return &preference, nil
}

func (pr *PreferenceRepository) CreateOrUpdate(preference *models.Preference) error {

	// Checking if the preference with the given username already exists
	existingPreference := &models.Preference{}
	result := pr.db.Where("username = ?", preference.Username).First(existingPreference)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Creating a new Preference if one does not exist,
			newPreference := models.Preference{
				Username:   preference.Username,
				Categories: preference.Categories,
			}
			if err := pr.db.Create(&newPreference).Error; err != nil {
				fmt.Printf("Error creating preference: %s\n", err.Error())
				return err
			}
		} else {
			// Other errors during the query
			fmt.Printf("Error querying preference: %s\n", result.Error.Error())
			return result.Error
		}
	} else {
		// updating Preference if it exists
		existingPreference.Categories = preference.Categories
		if err := pr.db.Save(existingPreference).Error; err != nil {
			fmt.Printf("Error updating preference: %s\n", err.Error())
			return err
		}
	}

	fmt.Printf("Preference created or updated successfully: %+v\n", preference)
	return nil
}
