package main

import (
	"dormitory_management/internal/database/db"
	"dormitory_management/internal/database/redis"
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"
	"dormitory_management/internal/server"
	"dormitory_management/internal/utils"
	"dormitory_management/pkg"
	"log"
)

func main() {
	pkg.LoadEnv()

	dbClient := db.NewDBClient()
	redisClient := redis.NewRedisClient()

	util := utils.NewUtil(dbClient, redisClient)
	handler := handlers.NewHandler(dbClient, redisClient, util)

	middleware := middleware.NewMiddleware(dbClient, redisClient, util)
	router := server.NewRouter(middleware, handler)

	log.Println("Server is running on port 8080")

	router.Run()
}
