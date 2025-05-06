package middleware

import (
	"dormitory_management/internal/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Middleware struct {
	dbClient    *gorm.DB
	redisClient *redis.Client
	util        *utils.Util
}

func NewMiddleware(dbClient *gorm.DB, redisClient *redis.Client, util *utils.Util) *Middleware {
	return &Middleware{dbClient: dbClient, redisClient: redisClient, util: util}
}
