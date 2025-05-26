package config

import (
	"log"
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
	Email    EmailConfig
	SMS      SMSConfig
}

// ServerConfig holds all the server-related configuration
type ServerConfig struct {
	Port      string
	Mode      string
	SecretKey string
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

type EmailConfig struct {
	FromEmail         string
	FromEmailPassword string
	FromEmailSMTP     string
	SMTP_ADDR         string
}

type SMSConfig struct {
	StringeeUrl          string
	StringeePhoneNumber  string
	StringeeISS          string
	StringeeTokenExpires int
	StringeeTokenSecret  string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Server config
	serverPort := getEnv("SERVER_PORT", "8080")
	serverMode := getEnv("SERVER_MODE", "debug")
	serverSecretKey := getEnv("SERVER_SECRET_KEY", "your-secret-key")

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
	redisDB := getEnvAsInt("REDIS_DB", 0)

	// JWT config
	jwtAccessSecret := getEnv("JWT_ACCESS_SECRET", "your-secret-key")
	jwtRefreshSecret := getEnv("JWT_REFRESH_SECRET", "your-secret-key")
	jwtAccessExpiresIn := getEnvAsInt("JWT_ACCESS_EXPIRES_IN", 3600)
	jwtRefreshExpiresIn := getEnvAsInt("JWT_REFRESH_EXPIRES_IN", 604800)

	// Log level
	logLevel := getEnv("LOG_LEVEL", "info")

	// Email config
	fromEmail := getEnv("FROM_EMAIL", "example@gmail.com")
	fromEmailPassword := getEnv("FROM_EMAIL_PASSWORD", "12345678")
	fromEmailSMTP := getEnv("FROM_EMAIL_SMTP", "smtp.gmail.com")
	smtp_ADDR := getEnv("SMTP_ADDR", "smtp.gmail.com:587")

	// SMS config
	stringeeUrl := getEnv("STRINGEE_URL", "https://api.stringee.com/v1/call2/callout")
	stringeePhoneNumber := getEnv("STRINGEE_PHONE_NUMBER", "+1234567890")
	stringeeISS := getEnv("STRINGEE_ISS", "+1234567890")
	stringeeTokenExpires := getEnvAsInt("STRINGEE_TOKEN_EXPIRES_IN", 3600)
	stringeeTokenSecret := getEnv("STRINGEE_TOKEN_SECRET", "your-secret-key")

	return &Config{
		Server: ServerConfig{
			Port:      serverPort,
			Mode:      serverMode,
			SecretKey: serverSecretKey,
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
		Email: EmailConfig{
			FromEmail:         fromEmail,
			FromEmailPassword: fromEmailPassword,
			FromEmailSMTP:     fromEmailSMTP,
			SMTP_ADDR:         smtp_ADDR,
		},
		SMS: SMSConfig{
			StringeeUrl:          stringeeUrl,
			StringeePhoneNumber:  stringeePhoneNumber,
			StringeeISS:          stringeeISS,
			StringeeTokenExpires: stringeeTokenExpires,
			StringeeTokenSecret:  stringeeTokenSecret,
		},
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

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
