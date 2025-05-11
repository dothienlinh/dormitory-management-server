package middleware

import (
	"dormitory_management/internal/services"
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
	response    *services.APIResponse
}

func NewMiddleware(dbClient *gorm.DB, redisClient *redis.Client, util *utils.Util, logger *zap.Logger, response *services.APIResponse) *Middleware {
	return &Middleware{dbClient: dbClient, redisClient: redisClient, util: util, logger: logger, response: response}
}
