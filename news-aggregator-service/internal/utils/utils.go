// utils.go
package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"news-aggregator-service/internal/models"
	"strconv"
	"time"
)

var jwtSecret []byte
var db *gorm.DB
var DefaultPreferences []string
var SpecialKey string
var ctx = context.Background()
var redisClient *redis.Client

// Representing the configuration for the News API
type NewsAPIConfig struct {
	URL        string   `mapstructure:"url"`
	Key        string   `mapstructure:"key"`
	Categories []string `mapstructure:"categories"`
}

// Representing the configuration for the Guardian News API
type GuardianNewsAPIConfig struct {
	URL      string   `mapstructure:"url"`
	Key      string   `mapstructure:"key"`
	Sections []string `mapstructure:"sections"`
}

var NewsAPI NewsAPIConfig
var GuardianNewsAPI GuardianNewsAPIConfig

// Initializing configuration
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

	// Reading News API config
	if err := viper.UnmarshalKey("news_api", &NewsAPI); err != nil {
		log.Fatalf("Error reading News API config: %v", err)
	}

	// Reading Guardian News API config
	if err := viper.UnmarshalKey("guardian_news_api", &GuardianNewsAPI); err != nil {
		log.Fatalf("Error reading Guardian News API config: %v", err)
	}
}

// loadConfig loads configuration from Viper
func loadConfig() {
	// Implement the function if needed
	// Example: Read configuration parameters from Viper
}

// DBConfig
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbname"`
	Password string `mapstructure:"password"`
}

// Retrieving the database configuration based on the environment
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

// Initializing the database connection
func InitDB() *gorm.DB {
	loadConfig()

	if db != nil {
		return db
	}

	dbConfig := getDBConfig()

	// Constructing connection string based on the retrieved configuration
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.DBName, dbConfig.Password)

	// Initializing database connection
	database, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// Enabling Logger, showing detailed log
	database.Logger.LogMode(logger.Info)

	// AutoMigrating models
	err = database.AutoMigrate(&models.NewsContent{})
	if err != nil {
		return nil
	}

	db = database
	return db
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

var Redis RedisConfig

// Initializing Redis client
func InitRedis() {
	InitConfig()

	Redis = RedisConfig{
		Address:  viper.GetString("redis.address"),
		Username: viper.GetString("redis.username"),
		Password: viper.GetString("redis.password"),
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     Redis.Address,
		Username: Redis.Username,
		Password: Redis.Password,
		DB:       0, // Use default DB
	})

	// Check connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func InitSpecialKey() {
	SpecialKey = viper.GetString("special_key")
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

func CalculateTotalPages(totalItems, pageSize int64) int {
	if totalItems <= 0 || pageSize <= 0 {
		return 0
	}

	totalPages := totalItems / pageSize
	if totalItems%pageSize > 0 {
		totalPages++
	}

	return int(totalPages)
}

func ParsePaginationParameters(c *gin.Context) (int, int, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		return 0, 0, errors.New("invalid page parameter")
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		return 0, 0, errors.New("invalid pageSize parameter")
	}

	return page, pageSize, nil
}

// ParseTimeParameter parses and validates a time parameter from the request.
func ParseTimeParameter(c *gin.Context, paramName string) (time.Time, error) {
	timeStr := c.Query(paramName)
	if timeStr == "" {
		return time.Time{}, nil
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, errors.New("invalid " + paramName + " parameter")
	}

	return parsedTime, nil
}

func ConvertToJSON(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
