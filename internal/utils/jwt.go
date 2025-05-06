package utils

import (
	"dormitory_management/internal/types"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Claims struct {
	UserID  uint   `json:"user_id"`
	TokenID string `json:"token_id"` // Unique ID for each token for blacklisting
	jwt.RegisteredClaims
}

var jwtAccessTokenSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))
var jwtRefreshTokenSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))

// 1 day
var jwtAccessTokenExpiration = time.Hour * 24 * 1

// 7 days
var jwtRefreshTokenExpiration = time.Hour * 24 * 7

func (u *Util) GenerateAccessToken(userID uint) (string, error) {
	now := time.Now()
	// Generate a unique token ID
	tokenID := uuid.New().String()

	claims := &Claims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtAccessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtAccessTokenSecret)
	if err != nil {
		return "", err
	}

	// Store token
	if err := u.StoreToken(userID, tokenID, types.AccessToken, jwtAccessTokenExpiration); err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}

	return tokenString, nil
}

func (u *Util) GenerateRefreshToken(userID uint) (string, error) {
	now := time.Now()
	// Generate a unique token ID
	tokenID := uuid.New().String()

	claims := &Claims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtRefreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtRefreshTokenSecret)
	if err != nil {
		return "", err
	}

	// Store token
	if err := u.StoreToken(userID, tokenID, types.RefreshToken, jwtRefreshTokenExpiration); err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}

	return tokenString, nil
}

func (u *Util) ValidateAccessToken(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtAccessTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if !u.ValidateToken(claims.UserID, claims.TokenID, types.AccessToken) {
		return nil, fmt.Errorf("token has been invalidated")
	}

	return claims, nil
}

func (u *Util) ValidateRefreshToken(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtRefreshTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if !u.ValidateToken(claims.UserID, claims.TokenID, types.RefreshToken) {
		return nil, fmt.Errorf("token has been invalidated")
	}

	return claims, nil
}

func (u *Util) InvalidateUserTokens(redisClient *redis.Client, userID uint) error {
	return u.InvalidateAllTokens(userID)
}

func (u *Util) StoreToken(userID uint, tokenID string, tokenType types.TokenType, expiration time.Duration) error {
	key := generateTokenKey(userID, tokenType)
	return u.redisClient.Set(u.ctx, key, tokenID, expiration).Err()
}

func (u *Util) ValidateToken(userID uint, tokenID string, tokenType types.TokenType) bool {
	key := generateTokenKey(userID, tokenType)
	storedToken, err := u.redisClient.Get(u.ctx, key).Result()
	if err != nil {
		return false
	}
	return storedToken == tokenID
}

func (u *Util) InvalidateToken(userID uint, tokenType types.TokenType) error {
	key := generateTokenKey(userID, tokenType)
	return u.redisClient.Del(u.ctx, key).Err()
}

func (u *Util) InvalidateAllTokens(userID uint) error {
	accessKey := generateAccessTokenKey(userID)
	refreshKey := generateRefreshTokenKey(userID)
	return u.redisClient.Del(u.ctx, accessKey, refreshKey).Err()
}

func generateTokenKey(userID uint, tokenType types.TokenType) string {
	switch tokenType {
	case types.AccessToken:
		return generateAccessTokenKey(userID)
	case types.RefreshToken:
		return generateRefreshTokenKey(userID)
	default:
		return ""
	}
}

func generateAccessTokenKey(userID uint) string {
	return fmt.Sprintf("user:access_token:%d", userID)
}

func generateRefreshTokenKey(userID uint) string {
	return fmt.Sprintf("user:refresh_token:%d", userID)
}
