package middleware

import (
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/helper"
	"dormitory_management/pkg/logger"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Claims is the custom JWT claims
type Claims struct {
	UserID       uint   `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TokenVersion string `json:"token_version"`
	jwt.RegisteredClaims
}

// Middleware contains all HTTP middleware handlers
type Middleware struct {
	repos  repository.Repositories
	logger logger.Logger
	config *config.Config
}

// NewMiddleware creates a new middleware instance
func NewMiddleware(repos repository.Repositories, logger logger.Logger) *Middleware {
	return &Middleware{
		repos:  repos,
		logger: logger,
		config: config.LoadConfig(),
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware logs each request
func (m *Middleware) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// After request
		latency := time.Since(start)
		status := c.Writer.Status()

		m.logger.Info("Request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// AuthMiddleware authenticates requests
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp := response.Unauthorized("Authorization header is required")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		// Check if the header has the "Bearer " prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			resp := response.Unauthorized("Invalid authorization format")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		decryptToken, err := helper.Decrypt(tokenString, m.config.Server.SecretKey)
		if err != nil {
			m.logger.Error("Failed to decrypt token", zap.Error(err))
			resp := response.Unauthorized(err.Error())
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(decryptToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.config.JWT.AccessSecret), nil
		})
		if err != nil {
			m.logger.Error("Failed to parse token", zap.Error(err))
			resp := response.Unauthorized(err.Error())
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			m.logger.Error("Invalid token")
			resp := response.Unauthorized("Invalid token")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		// Extract the claims from the token
		claims, ok := token.Claims.(*Claims)
		if !ok {
			m.logger.Error("Failed to extract claims from token")
			resp := response.Unauthorized("Failed to extract claims from token")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		tokenVersion, err := m.repos.Auth().CheckTokenVersion(c, entity.AccessToken, claims.UserID)
		if err != nil {
			m.logger.Error("Failed to check token version", zap.Error(err))
			resp := response.Unauthorized(err.Error())
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		if tokenVersion != claims.TokenVersion {
			m.logger.Error("Invalid token version", zap.Error(err))
			resp := response.Unauthorized("Invalid token version")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		// Get user
		user, err := m.repos.Auth().GetUserCache(c, claims.UserID)
		if errors.Is(err, redis.Nil) {
			user, err := m.repos.User().GetByID(c, claims.UserID)
			if err != nil {
				resp := response.Unauthorized(err.Error())
				c.JSON(resp.Status, resp.Response)
				c.Abort()
				return
			}

			// Set user cache
			err = m.repos.Auth().SetUserCache(c, user, m.config.JWT.AccessExpiresIn)
			if err != nil {
				m.logger.Error("Failed to set user cache", zap.Error(err))
				resp := response.InternalServerError(err.Error())
				c.JSON(resp.Status, resp.Response)
				c.Abort()
				return
			}

			// Set user in context
			m.setUserInContext(c, user)
			c.Next()
			return
		}

		if err != nil {
			resp := response.Unauthorized(err.Error())
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		// Set the user ID and role in the context
		m.setUserInContext(c, user)

		c.Next()
	}
}

func (m *Middleware) setUserInContext(c *gin.Context, user *entity.User) {
	c.Set("userID", user.ID)
	c.Set("userRole", user.Role)
	c.Set("userEmail", user.Email)
}

// AdminMiddleware ensures the user has admin role
func (m *Middleware) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This assumes AuthMiddleware has already been run
		role, exists := c.Get("userRole")
		if !exists {
			resp := response.Unauthorized("User role not found")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		// Check if user is admin
		if role != string(entity.UserRoleAdmin) {
			resp := response.Unauthorized("Admin access required")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		c.Next()
	}
}

// ErrorMiddleware handles errors globally
func (m *Middleware) ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only handle errors if there are errors to handle
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			m.logger.Error("API error", zap.Error(err))

			var resp response.StatusResponse

			// Determine the type of error and set the appropriate status code and message
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				resp = response.NotFound("Resource not found")
			default:
				resp = response.InternalServerError("Internal server error")
			}

			c.JSON(resp.Status, resp.Response)
		}
	}
}

func (m *Middleware) CustomRecovery() gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		resp := response.InternalServerError("Internal server error")
		c.JSON(resp.Status, resp.Response)
	}
}
