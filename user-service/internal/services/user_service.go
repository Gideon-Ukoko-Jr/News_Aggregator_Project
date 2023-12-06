package services

import (
	"fmt"
	"user-service/internal/models"
	"user-service/internal/repositories"
	"user-service/internal/utils"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

// NewUserService
func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) RegisterUser(user *models.User) error {
	// Hashing the user's password before saving it to the db
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	//Debug
	//fmt.Printf("Registration - Hashed Password: %s\n", hashedPassword)

	user.Password = hashedPassword

	return us.userRepository.CreateUser(user)
}

func (us *UserService) LoginUser(username, password string) (*models.User, error) {
	user, err := us.userRepository.GetUserByUsername(username)
	if err != nil {
		fmt.Printf("LoginUser - Error getting user by username: %v\n", err)
		return nil, err
	}

	if user == nil {
		fmt.Printf("LoginUser - User not found for username: %s\n", username)
		return nil, utils.ErrInvalidCredentials
	}

	// Validating the password
	if !utils.ValidatePassword(password, user.Password) {
		fmt.Printf("LoginUser - Password validation failed for username: %s\n", username)
		return nil, utils.ErrInvalidCredentials
	}

	return user, nil
}

func (us *UserService) GetUserByID(authUserID, requestedUserID uint) (*models.User, error) {
	if authUserID != requestedUserID {
		return nil, utils.ErrUnauthorizedAccess
	}

	return us.userRepository.GetUserByID(requestedUserID)
}
