package models

import "gorm.io/gorm"

// User model
type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null;unique" json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token,omitempty"`
}

// NewUserResponse
func NewUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    user.Token,
	}
}

// UserResponse
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
}

// RegisterRequest
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
