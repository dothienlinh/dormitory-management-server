package middleware

import (
	"dormitory_management/internal/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Middleware struct {
	dbClient    *gorm.DB
	redisClient *redis.Client
	util        *utils.Util
	logger      *zap.Logger
}

func NewMiddleware(dbClient *gorm.DB, redisClient *redis.Client, util *utils.Util, logger *zap.Logger) *Middleware {
	return &Middleware{dbClient: dbClient, redisClient: redisClient, util: util, logger: logger}
}
