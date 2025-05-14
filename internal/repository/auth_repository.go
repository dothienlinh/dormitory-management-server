package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/entity/helper"
	"dormitory_management/internal/domain/repository"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// authRepository implements the repository.AuthRepository interface
type authRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

// NewAuthRepository creates a new auth repository
func NewAuthRepository(db *gorm.DB, redisClient *redis.Client) repository.AuthRepository {
	return &authRepository{
		db:          db,
		redisClient: redisClient,
	}
}

// StoreToken stores a token in Redis
func (r *authRepository) StoreToken(ctx context.Context, userID uint, token string, expiresIn int) error {
	key := fmt.Sprintf("refresh_token:%s", token)
	return r.redisClient.Set(ctx, key, userID, time.Duration(expiresIn)*time.Second).Err()
}

// GetToken retrieves a token from Redis
func (r *authRepository) GetToken(ctx context.Context, token string) (uint, error) {
	key := fmt.Sprintf("refresh_token:%s", token)
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, errors.New("token not found or expired")
		}
		return 0, fmt.Errorf("failed to get token: %w", err)
	}

	userID, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return uint(userID), nil
}

// InvalidateToken invalidates a token in Redis
func (r *authRepository) InvalidateToken(ctx context.Context, userID uint, token string) error {
	key := fmt.Sprintf("refresh_token:%s", token)
	return r.redisClient.Del(ctx, key).Err()
}

// Register registers a new user
func (r *authRepository) Register(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}
	return nil
}

// Login authenticates a user and returns user data
func (r *authRepository) Login(ctx context.Context, email, password string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if !helper.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}
