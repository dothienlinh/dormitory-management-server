package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	LogLevel string
}

// ServerConfig holds all the server-related configuration
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig holds all the database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// RedisConfig holds all the Redis-related configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds all the JWT-related configuration
type JWTConfig struct {
	AccessSecret     string
	RefreshSecret    string
	AccessExpiresIn  int
	RefreshExpiresIn int
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	godotenv.Load()

	// Server config
	serverPort := getEnv("SERVER_PORT", "8080")
	serverMode := getEnv("SERVER_MODE", "debug")

	// Database config
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "dormitory")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")

	// Redis config
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	// JWT config
	jwtAccessSecret := getEnv("JWT_ACCESS_SECRET", "your-secret-key")
	jwtRefreshSecret := getEnv("JWT_REFRESH_SECRET", "your-secret-key")
	jwtAccessExpiresIn, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRES_IN", "3600"))
	jwtRefreshExpiresIn, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRES_IN", "604800"))

	// Log level
	logLevel := getEnv("LOG_LEVEL", "info")

	return &Config{
		Server: ServerConfig{
			Port: serverPort,
			Mode: serverMode,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
		Redis: RedisConfig{
			Host:     redisHost,
			Port:     redisPort,
			Password: redisPassword,
			DB:       redisDB,
		},
		JWT: JWTConfig{
			AccessSecret:     jwtAccessSecret,
			RefreshSecret:    jwtRefreshSecret,
			AccessExpiresIn:  jwtAccessExpiresIn,
			RefreshExpiresIn: jwtRefreshExpiresIn,
		},
		LogLevel: logLevel,
	}
}

// getEnv retrieves the value of the environment variable named by the key
// or returns the fallback value if the variable is not set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
