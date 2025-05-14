package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/helper"
	"errors"
	"fmt"
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

// CheckTokenVersion checks the token version in Redis
func (r *authRepository) CheckTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint) (string, error) {
	key := fmt.Sprintf("%s:%d", tokenType.String(), userID)
	return r.redisClient.Get(ctx, key).Result()

}

// SetTokenVersion sets the token version in Redis
func (r *authRepository) SetTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint, tokenVersion string, expiresIn int) error {
	key := fmt.Sprintf("%s:%d", tokenType.String(), userID)
	return r.redisClient.Set(ctx, key, tokenVersion, time.Duration(expiresIn)*time.Second).Err()
}

// InvalidateToken invalidates a token in Redis
func (r *authRepository) InvalidateToken(ctx context.Context, tokenType entity.TokenType, userID uint) error {
	key := fmt.Sprintf("%s:%d", tokenType.String(), userID)
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
