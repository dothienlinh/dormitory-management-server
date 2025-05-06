package handlers

import (
	"dormitory_management/internal/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	dbClient    *gorm.DB
	redisClient *redis.Client
	util        *utils.Util
}

func NewHandler(dbClient *gorm.DB, redisClient *redis.Client, util *utils.Util) *Handler {
	return &Handler{dbClient: dbClient, redisClient: redisClient, util: util}
}
