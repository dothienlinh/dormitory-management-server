package utils

import (
	"dormitory_management/internal/database/redis"
	"dormitory_management/internal/types"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	StoreToken(userID uint, tokenID string, tokenType types.TokenType, expiration time.Duration) error
	ValidateToken(userID uint, tokenID string, tokenType types.TokenType) bool
	InvalidateAllTokens(userID uint) error
}

var tokenManager TokenManager = redis.NewRedisTokenManager()

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

func GenerateAccessToken(userID uint) (string, error) {
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
	if tokenManager != nil {
		if err := tokenManager.StoreToken(userID, tokenID, types.AccessToken, jwtAccessTokenExpiration); err != nil {
			return "", fmt.Errorf("failed to store token: %w", err)
		}
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID uint) (string, error) {
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
	if tokenManager != nil {
		if err := tokenManager.StoreToken(userID, tokenID, types.RefreshToken, jwtRefreshTokenExpiration); err != nil {
			return "", fmt.Errorf("failed to store token: %w", err)
		}
	}

	return tokenString, nil
}

func ValidateAccessToken(token string) (*Claims, error) {
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
	if tokenManager != nil && !tokenManager.ValidateToken(claims.UserID, claims.TokenID, types.AccessToken) {
		return nil, fmt.Errorf("token has been invalidated")
	}

	return claims, nil
}

func ValidateRefreshToken(token string) (*Claims, error) {
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
	if tokenManager != nil && !tokenManager.ValidateToken(claims.UserID, claims.TokenID, types.RefreshToken) {
		return nil, fmt.Errorf("token has been invalidated")
	}

	return claims, nil
}

func InvalidateUserTokens(userID uint) error {
	if tokenManager == nil {
		return nil
	}
	return tokenManager.InvalidateAllTokens(userID)
}
