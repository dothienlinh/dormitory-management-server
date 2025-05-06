package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Util struct {
	dbClient    *gorm.DB
	redisClient *redis.Client
	ctx         context.Context
}

func NewUtil(dbClient *gorm.DB, redisClient *redis.Client) *Util {
	return &Util{dbClient: dbClient, redisClient: redisClient, ctx: context.Background()}
}
