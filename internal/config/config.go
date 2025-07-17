package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	LogLevel string
	Email    EmailConfig
	Client   ClientConfig
	PayOS    PayOS
}

type ServerConfig struct {
	Port      string
	Mode      string
	SecretKey string
}

type ClientConfig struct {
	ClientDomain string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

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

type PayOS struct {
	ClientID    string
	APIKey      string
	ChecksumKey string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	appEnv := getEnv("APP_ENV", "development")

	serverPort := getEnv("SERVER_PORT", "8080")
	serverMode := getEnv("SERVER_MODE", "debug")
	serverSecretKey := getEnv("SERVER_SECRET_KEY", "your-secret-key")

	clientDomain := getEnv("CLIENT_DOMAIN", "http://localhost:3000")

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "dormitory")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := getEnvAsInt("REDIS_DB", 0)

	jwtAccessSecret := getEnv("JWT_ACCESS_SECRET", "your-secret-key")
	jwtRefreshSecret := getEnv("JWT_REFRESH_SECRET", "your-secret-key")
	jwtAccessExpiresIn := getEnvAsInt("JWT_ACCESS_EXPIRES_IN", 3600)
	jwtRefreshExpiresIn := getEnvAsInt("JWT_REFRESH_EXPIRES_IN", 604800)

	logLevel := getEnv("LOG_LEVEL", "info")

	fromEmail := getEnv("FROM_EMAIL", "example@gmail.com")
	fromEmailPassword := getEnv("FROM_EMAIL_PASSWORD", "12345678")
	fromEmailSMTP := getEnv("FROM_EMAIL_SMTP", "smtp.gmail.com")
	smtp_ADDR := getEnv("SMTP_ADDR", "smtp.gmail.com:587")

	clientID := getEnv("PAYOS_CLIENT_ID", "")
	apiKey := getEnv("PAYOS_API_KEY", "")
	checksumKey := getEnv("PAYOS_CHECKSUM_KEY", "")

	return &Config{
		AppEnv: appEnv,
		Server: ServerConfig{
			Port:      serverPort,
			Mode:      serverMode,
			SecretKey: serverSecretKey,
		},
		Client: ClientConfig{
			ClientDomain: clientDomain,
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
		PayOS: PayOS{
			ClientID:    clientID,
			APIKey:      apiKey,
			ChecksumKey: checksumKey,
		},
	}
}

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
