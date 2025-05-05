package redis

import (
	"context"
	"dormitory_management/internal/types"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenManager struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisTokenManager() *RedisTokenManager {
	return &RedisTokenManager{
		client: RedisClient,
		ctx:    context.Background(),
	}
}

func (m *RedisTokenManager) StoreToken(userID uint, tokenID string, tokenType types.TokenType, expiration time.Duration) error {
	key := m.generateTokenKey(userID, tokenType)
	return m.client.Set(m.ctx, key, tokenID, expiration).Err()
}

func (m *RedisTokenManager) ValidateToken(userID uint, tokenID string, tokenType types.TokenType) bool {
	key := m.generateTokenKey(userID, tokenType)
	storedToken, err := m.client.Get(m.ctx, key).Result()
	if err != nil {
		return false
	}
	return storedToken == tokenID
}

func (m *RedisTokenManager) InvalidateToken(userID uint, tokenType types.TokenType) error {
	key := m.generateTokenKey(userID, tokenType)
	return m.client.Del(m.ctx, key).Err()
}

func (m *RedisTokenManager) InvalidateAllTokens(userID uint) error {
	accessKey := m.generateAccessTokenKey(userID)
	refreshKey := m.generateRefreshTokenKey(userID)
	return m.client.Del(m.ctx, accessKey, refreshKey).Err()
}

func (m *RedisTokenManager) generateTokenKey(userID uint, tokenType types.TokenType) string {
	switch tokenType {
	case types.AccessToken:
		return m.generateAccessTokenKey(userID)
	case types.RefreshToken:
		return m.generateRefreshTokenKey(userID)
	default:
		return ""
	}
}

func (m *RedisTokenManager) generateAccessTokenKey(userID uint) string {
	return fmt.Sprintf("user:access_token:%d", userID)
}

func (m *RedisTokenManager) generateRefreshTokenKey(userID uint) string {
	return fmt.Sprintf("user:refresh_token:%d", userID)
}
