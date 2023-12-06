package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
	"user-service/internal/handlers"
	"user-service/internal/middleware"
	"user-service/internal/repositories"
	"user-service/internal/services"
	"user-service/internal/utils"
)

const apiPrefix = "/api/v1/user"

func main() {
	// Initializing configuration
	utils.InitConfig()

	// Initializing database
	db := utils.InitDB()
	sqlDB, err := db.DB()
	if err != nil {
		// Handle error
		panic(err)
	}
	defer sqlDB.Close()

	// Initializing Gin router
	router := gin.Default()

	// Initializing JWT secret key
	utils.InitJWTSecret()

	// Applying middleware for JWT authentication
	router.Use(middleware.AuthMiddleware())

	// Initializing repositories and services
	userRepository := repositories.NewUserRepository(db)
	preferenceRepository := repositories.NewPreferenceRepository(db)
	userService := services.NewUserService(userRepository)
	preferenceService := services.NewPreferenceService(preferenceRepository)

	// Initializing handlers with services
	userHandler := handlers.NewUserHandler(userService, preferenceService)

	// Routes
	router.POST(apiPrefix+"/register", userHandler.Register)
	router.POST(apiPrefix+"/login", userHandler.Login)

	authGroup := router.Group(apiPrefix)
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/user-details/:id", userHandler.GetUserByID)
		authGroup.GET("/user-news-preferences", userHandler.GetUserPreferences)
		authGroup.PUT("/update-news-preferences", userHandler.SetPreferences)
	}

	fmt.Println("Server is running on :8080")
	router.Run(":8080")
}
