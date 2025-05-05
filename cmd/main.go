package main

import (
	"dormitory_management/internal/database/db"
	"dormitory_management/internal/database/redis"
	"dormitory_management/internal/server"
	"dormitory_management/pkg"
)

func main() {
	pkg.LoadEnv()

	db.ConnectDatabase()
	redis.ConnectRedis()

	server.SetupRouter()
}
