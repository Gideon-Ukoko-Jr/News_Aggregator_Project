package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"net/http"
	"strings"
	"user-service/internal/models"
	"user-service/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("AuthMiddleware: Checking JWT token")

		// Extracting token from the Authorization header
		tokenString := extractTokenFromHeader(c)
		if tokenString == "" {
			// No token provided, continue without checking
			c.Next()
			return
		}

		// Validating and parsing the JWT token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			// Token validation failed, handle accordingly (e.g., return an error)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Extracting user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			// Handling User ID not found in claims
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Setting authenticated user in the context
		c.Set("user", &models.User{Model: gorm.Model{ID: uint(userID)}})

		// Continuing with the request
		c.Next()
	}
}

// Extracting the JWT token from the Authorization header
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Checking if the Authorization header has the "Bearer" prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
