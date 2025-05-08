package main

import (
	"dormitory_management/internal/database/db"
	"dormitory_management/internal/database/redis"
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"
	"dormitory_management/internal/server"
	"dormitory_management/internal/utils"
	"dormitory_management/pkg"
)

func main() {
	pkg.LoadEnv()
	logger := pkg.ConfigureLogger()

	dbClient := db.NewDBClient()
	redisClient := redis.NewRedisClient()

	util := utils.NewUtil(dbClient, redisClient)
	handler := handlers.NewHandler(dbClient, redisClient, util, logger)

	middleware := middleware.NewMiddleware(dbClient, redisClient, util, logger)
	router := server.NewRouter(middleware, handler)

	logger.Info("Server is running on port 8080")

	defer logger.Sync()

	router.Run()
}
