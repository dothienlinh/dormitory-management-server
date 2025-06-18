package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/helper"
	"dormitory_management/internal/infra/cache"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// authRepository implements the repository.AuthRepository interface
type authRepository struct {
	db    *gorm.DB
	redis *cache.RedisClient
}

// NewAuthRepository creates a new auth repository
func NewAuthRepository(db *gorm.DB, redisClient *cache.RedisClient) repository.AuthRepository {
	return &authRepository{
		db:    db,
		redis: redisClient,
	}
}

// CheckTokenVersion checks the token version in Redis
func (r *authRepository) CheckTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint64) (string, error) {
	key := fmt.Sprintf("%s:%d", tokenType.String(), userID)
	return r.redis.Get(ctx, key)

}

// SetTokenVersion sets the token version in Redis
func (r *authRepository) SetCacheTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint64, tokenVersion string, expiresIn int) error {
	key := fmt.Sprintf("%s:%d", tokenType.String(), userID)
	return r.redis.Set(ctx, key, tokenVersion, expiresIn)
}

// SetUserCache sets user data in Redis
func (r *authRepository) SetUserCache(ctx context.Context, user *entity.User, expiresIn int) error {
	key := fmt.Sprintf("user:%d", user.ID)
	userJson, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}
	return r.redis.Set(ctx, key, string(userJson), expiresIn)
}

// GetUserCache retrieves user data from Redis
func (r *authRepository) GetUserCache(ctx context.Context, user *entity.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	userJson, err := r.redis.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(userJson), user); err != nil {
		return err
	}
	return nil
}

// DeleteUserCache deletes user data from Redis
func (r *authRepository) DeleteUserCache(ctx context.Context, userID uint64) error {
	key := fmt.Sprintf("user:%d", userID)
	return r.redis.Del(ctx, key)
}

// InvalidateToken invalidates a token in Redis
func (r *authRepository) InvalidateToken(ctx context.Context, tokenType entity.TokenType, userID uint64) error {
	key := fmt.Sprintf("%s:%d", tokenType.String(), userID)
	return r.redis.Del(ctx, key)
}

// Register registers a new user
func (r *authRepository) Register(ctx context.Context, user *entity.User, otpCode *entity.OtpCode) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Table(user.TableName()).Create(user).Error; err != nil {
			return fmt.Errorf("failed to register user: %w", err)
		}

		if err := tx.WithContext(ctx).Table(otpCode.TableName()).Create(otpCode).Error; err != nil {
			return err
		}

		return nil
	})

}

// Login authenticates a user and returns user data
func (r *authRepository) Login(ctx context.Context, user *entity.User) error {
	password := user.Password
	if err := r.db.WithContext(ctx).Table(user.TableName()).Where("email = ?", user.Email).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid email or password")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsVerify {
		return errors.New("unverified user")
	}

	// Verify password
	if !helper.CheckPassword(password, user.Password) {
		return errors.New("invalid email or password")
	}

	return nil
}

func (r *authRepository) VerifyAccount(ctx context.Context, otpCode *entity.OtpCode, user *entity.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		useOtpCode := entity.UseOtpCode{IsUsed: true, VerifiedAt: time.Now().Format(time.RFC3339)}
		if err := tx.WithContext(ctx).Table(otpCode.TableName()).Where(otpCode).Updates(useOtpCode).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Table(user.TableName()).Where(user).Updates(entity.User{IsVerify: true}).Error; err != nil {
			return err
		}

		return nil
	})
}
