package services

import (
	"fmt"
	"user-service/internal/models"
	"user-service/internal/repositories"
	"user-service/internal/utils"
)

type PreferenceService struct {
	preferenceRepository *repositories.PreferenceRepository
}

func NewPreferenceService(preferenceRepository *repositories.PreferenceRepository) *PreferenceService {
	return &PreferenceService{preferenceRepository: preferenceRepository}
}

func (ps *PreferenceService) GetUserPreferences(username string) ([]string, error) {
	preference, err := ps.preferenceRepository.FindByUsername(username)
	//if err != nil {
	//	return nil, err
	//}

	if preference == nil && err != nil {
		// If no preference record found, create a new one with default preferences
		newPreference := &models.Preference{
			Username:   username,
			Categories: utils.GetDefaultPreferences(),
		}

		// Try to create the new preference record
		if err := ps.preferenceRepository.CreateOrUpdate(newPreference); err != nil {
			fmt.Printf("Error creating new preference record: %v\n", err)
			return nil, err
		}

		// Return the default preferences
		return newPreference.Categories, nil
	}

	return preference.Categories, nil
}

func (ps *PreferenceService) SetUserPreferences(username string, preferences *models.Preference) error {
	preference, err := ps.preferenceRepository.FindByUsername(username)
	if err != nil {
		return err
	}

	if preference == nil {
		// If no preference record found, create a new one with default preferences
		preference = &models.Preference{
			Username:   username,
			Categories: utils.GetDefaultPreferences(),
		}
	}

	preference.Categories = preferences.Categories
	return ps.preferenceRepository.CreateOrUpdate(preference)
}
