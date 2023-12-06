package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
	"user-service/internal/models"
)

var jwtSecret []byte
var db *gorm.DB
var DefaultPreferences []string

// InitConfig initializes configuration
func InitConfig() {
	viper.SetConfigName("config") // name of your config file (without extension)
	viper.AddConfigPath("config") // path to look for the config file in
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Set default JWT secret if not provided in the config
	viper.SetDefault("jwt_secret", "defaultsecretKey$23")

	// Reading default preferences
	DefaultPreferences = viper.GetStringSlice("default_preferences")
}

// loadConfig loads configuration from Viper
func loadConfig() {
	// Implement the function if needed
	// Example: Read configuration parameters from Viper
}

// DBConfig represents the database configuration
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbname"`
	Password string `mapstructure:"password"`
}

// getDBConfig retrieves the database configuration based on the environment
func getDBConfig() DBConfig {
	useLocal := viper.GetBool("database.useLocal")

	if useLocal {
		return DBConfig{
			Host:     viper.GetString("database.local.host"),
			Port:     viper.GetInt("database.local.port"),
			User:     viper.GetString("database.local.user"),
			DBName:   viper.GetString("database.local.dbname"),
			Password: viper.GetString("database.local.password"),
		}
	}

	return DBConfig{
		Host:     viper.GetString("database.prod.host"),
		Port:     viper.GetInt("database.prod.port"),
		User:     viper.GetString("database.prod.user"),
		DBName:   viper.GetString("database.prod.dbname"),
		Password: viper.GetString("database.prod.password"),
	}
}

// InitDB initializes the database connection
func InitDB() *gorm.DB {
	loadConfig()

	if db != nil {
		return db
	}

	dbConfig := getDBConfig()

	// Construct connection string based on the retrieved configuration
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.DBName, dbConfig.Password)

	// Initialize database connection
	database, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// Enable Logger, show detailed log
	database.Logger.LogMode(logger.Info)

	// AutoMigrate models
	err = database.AutoMigrate(&models.User{}, &models.Preference{})
	if err != nil {
		return nil
	}

	db = database
	return db
}

// Initializing JWT secret key
func InitJWTSecret() {
	InitConfig()

	// Reading JWT secret from the config file
	jwtSecret = []byte(viper.GetString("jwt_secret"))
	//fmt.Println("JWT Secret:", string(jwtSecret))
}

// Generating a JWT token for a user
func GenerateJWT(userID uint) (string, error) {

	expirationTime := time.Now().Add(3 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}

	return signedToken, nil
}

// Validating a JWT token
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	return string(hashedPassword), nil
}

// Validating a password against a hashed password
func ValidatePassword(password, hashedPassword string) bool {
	fmt.Printf("Comparing Password: %s with Hashed Password: %s\n", password, hashedPassword)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Printf("Password validation failed: %v\n", err)
	}
	return err == nil
}

func GetDefaultPreferences() []string {
	return DefaultPreferences
}

func ValidateCategories(categories []string) error {
	defaultPreferences := GetDefaultPreferences()

	for _, category := range categories {
		if !contains(defaultPreferences, category) {
			return errors.New("invalid category: " + category)
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

var ErrInvalidCredentials = errors.New("invalid login credentials")
var ErrDuplicateEmail = errors.New("email is already registered")
var ErrUnauthorizedAccess = errors.New("unauthorized access")
