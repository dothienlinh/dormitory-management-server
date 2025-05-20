package cache

import (
	"context"
	"dormitory_management/internal/config"
	"dormitory_management/pkg/logger"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClientInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
}

type RedisClient struct {
	*redis.Client
	logger logger.Logger
}

// NewRedisClient creates a new Redis client connection
func NewRedisClient(cfg config.RedisConfig, logger logger.Logger, ctx context.Context) (*RedisClient, error) {
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	redisClient := &RedisClient{Client: client, logger: logger}

	if err := redisClient.Client.Ping(ctx).Err(); err != nil {
		logger.Error("failed to ping Redis server", zap.Error(err))
		return nil, err
	}

	return redisClient, nil
}

// Set sets a key-value pair in Redis with an expiration time
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	ttl := time.Duration(expiration) * time.Second
	if err := r.Client.Set(ctx, key, value, ttl).Err(); err != nil {
		r.logger.Error("failed to set key in Redis",
			zap.String("key", key),
			zap.Any("value", value),
			zap.Error(err))
		return err
	}

	r.logger.Info("set key in Redis",
		zap.String("key", key),
		zap.Any("value", value),
		zap.Duration("expiration", ttl))

	return nil
}

// Get retrieves a value from Redis by key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		r.logger.Error("failed to get key from Redis", zap.String("key", key), zap.Error(err))
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

// Del deletes a key from Redis
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	if err := r.Client.Del(ctx, keys...).Err(); err != nil {
		r.logger.Error("failed to delete key from Redis", zap.Strings("keys", keys), zap.Error(err))
		return err
	}
	return nil
}
