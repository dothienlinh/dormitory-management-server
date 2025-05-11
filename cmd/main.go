package main

import (
	"dormitory_management/internal/database/db"
	"dormitory_management/internal/database/redis"
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"
	"dormitory_management/internal/server"
	"dormitory_management/internal/services"
	"dormitory_management/internal/utils"
	"dormitory_management/pkg"
)

func main() {
	pkg.LoadEnv()

	logger := pkg.ConfigureLogger()

	dbClient := db.NewDBClient()

	redisClient := redis.NewRedisClient()

	util := utils.NewUtil(dbClient, redisClient)

	middleware := middleware.NewMiddleware(dbClient, redisClient, util, logger)

	service := services.NewService(dbClient, redisClient, logger, util)

	handler := handlers.NewHandler(service, logger)

	router := server.NewRouter(middleware, handler)

	defer logger.Sync()

	logger.Info("Server is running on port 8080")

	router.Run()
}
