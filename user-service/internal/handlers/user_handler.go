package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-service/internal/models"
	"user-service/internal/services"
	"user-service/internal/utils"
)

type UserHandler struct {
	userService       *services.UserService
	preferenceService *services.PreferenceService
}

// NewUserHandler
func NewUserHandler(userService *services.UserService, preferenceService *services.PreferenceService) *UserHandler {
	return &UserHandler{userService: userService, preferenceService: preferenceService}
}

// @Summary Register a new user
// @Description Register a new user with the given credentials
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.RegisterRequest true "User registration details"
// @Success 200 {object} models.UserResponse
// @Router /register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	var registerRequest models.RegisterRequest

	// Binding request body to User struct
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//// Debug - Logging the username and password for testing
	//fmt.Printf("Registration - Username: %s, Password: %s\n", user.Username, user.Password)

	if registerRequest.Username == "" || registerRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})
		return
	}

	// Registering user
	user := models.User{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
	}

	if err := uh.userService.RegisterUser(&user); err != nil {
		if err == utils.ErrDuplicateEmail {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate email address"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Returning User Response
	userResponse := models.NewUserResponse(&user)

	c.JSON(http.StatusOK, userResponse)
}

// @Summary Log in a user
// @Description Log in a user with credentials
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "User Login details"
// @Success 200 {object} models.UserResponse
// @Router /login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	var loginRequest models.LoginRequest

	// Binding request body to LoginRequest struct
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//// Debug - Logging the username and password for testing
	//fmt.Printf("Login - Username: %s, Password: %s\n", loginRequest.Username, loginRequest.Password)

	if loginRequest.Username == "" || loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})
		return
	}

	// Authenticating user
	user, err := uh.userService.LoginUser(loginRequest.Username, loginRequest.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generating JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Returning user response with token
	userResponse := models.NewUserResponse(user)
	userResponse.Token = token

	c.JSON(http.StatusOK, userResponse)
}

// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]string "error": "Invalid user ID"
// @Failure 401 {object} map[string]string "error": "Unauthorized"
// @Failure 403 {object} map[string]string "error": "Forbidden"
// @Failure 404 {object} map[string]string "error": "User not found"
// @Router /users/{id} [get]
func (uh *UserHandler) GetUserByID(c *gin.Context) {
	// Getting user ID from path parameters
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Getting authenticated user from auth middleware
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := uh.userService.GetUserByID(authUser.(*models.User).ID, uint(userID))
	if err != nil {
		if err == utils.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Creating user response without sensitive information
	userResponse := models.NewUserResponse(user)

	c.JSON(http.StatusOK, userResponse)
}

func (uh *UserHandler) GetUserPreferences(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetching user to get the username
	user, err := uh.userService.GetUserByID(authUser.(*models.User).ID, authUser.(*models.User).ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetching preferences for the user using the correct username
	preferences, err := uh.preferenceService.GetUserPreferences(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preferences": preferences})
}

func (uh *UserHandler) SetPreferences(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Fetching user to get the username
	user, err := uh.userService.GetUserByID(authUser.(*models.User).ID, authUser.(*models.User).ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Binding request body to PreferencesRequest struct
	var preferencesRequest models.PreferencesRequest
	if err := c.ShouldBindJSON(&preferencesRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validating categories
	if err := utils.ValidateCategories(preferencesRequest.Categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//fmt.Printf(user.Username)
	// Setting preferences for the user
	preferences := models.Preference{
		Username:   user.Username,
		Categories: preferencesRequest.Categories,
	}
	err = uh.preferenceService.SetUserPreferences(user.Username, &preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences set successfully"})
}
