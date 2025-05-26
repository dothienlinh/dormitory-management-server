package main

import (
	"context"
	"dormitory_management/internal/config"
	"dormitory_management/internal/delivery/http"
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"
	"dormitory_management/internal/delivery/mq"
	"dormitory_management/internal/infra/cache"
	"dormitory_management/internal/infra/database"
	"dormitory_management/internal/repository"
	"dormitory_management/internal/usecase"
	"dormitory_management/pkg/logger"
)

func main() {
	context := context.Background()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	log := logger.NewLogger(cfg.LogLevel)

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Initialize Redis client
	rdb, err := cache.NewRedisClient(cfg.Redis, log, context)
	if err != nil {
		log.Fatal("Failed to connect to Redis", err)
	}

	// Initialize Asynq client
	asynqClient := mq.NewClient(cfg)

	// Initialize repositories
	repos := repository.NewRepositories(db, rdb)

	// Initialize use cases
	useCases := usecase.NewUseCases(repos, log, asynqClient.Client())

	// Initialize middleware
	mw := middleware.NewMiddleware(repos, log)

	// Initialize HTTP handlers
	handlers := handler.NewHandlers(useCases, log)

	// Setup and run the server
	server := http.NewServer(cfg, handlers, mw)

	log.Info("Server is running on port " + cfg.Server.Port)
	if err := server.Run(); err != nil {
		log.Fatal("Server failed to start", err)
	}

	defer func() {
		if err := log.Sync(); err != nil {
			log.Fatal("Error syncing logger", err)
		}
	}()

	defer func() {
		if err := asynqClient.Close(); err != nil {
			log.Fatal("Error closing Asynq client", err)
		}
	}()
}
