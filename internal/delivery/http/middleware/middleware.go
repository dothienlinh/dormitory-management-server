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
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Claims struct {
	UserID       uint64 `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TokenVersion string `json:"token_version"`
	jwt.RegisteredClaims
}

type Middleware struct {
	repos  repository.Repositories
	logger logger.Logger
	config *config.Config
}

func NewMiddleware(repos repository.Repositories, logger logger.Logger) *Middleware {
	return &Middleware{
		repos:  repos,
		logger: logger,
		config: config.LoadConfig(),
	}
}

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

func (m *Middleware) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

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

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp := response.Unauthorized("Authorization header is required")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			resp := response.Unauthorized("Invalid authorization format")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		decryptToken, err := helper.Decrypt(tokenString, m.config.Server.SecretKey)
		if err != nil {
			m.logger.Error("Failed to decrypt token", zap.Error(err))
			resp := response.Unauthorized(err.Error())
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(decryptToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

		if !token.Valid {
			m.logger.Error("Invalid token")
			resp := response.Unauthorized("Invalid token")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

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

		user := &entity.User{
			Base: entity.Base{
				ID: claims.UserID,
			},
		}
		if err := m.repos.Auth().GetUserCache(c, user); errors.Is(err, redis.Nil) {
			if err := m.repos.User().GetByID(c, user); err != nil {
				resp := response.Unauthorized(err.Error())
				c.JSON(resp.Status, resp.Response)
				c.Abort()
				return
			}

			err = m.repos.Auth().SetUserCache(c, user, m.config.JWT.AccessExpiresIn)
			if err != nil {
				m.logger.Error("Failed to set user cache", zap.Error(err))
				resp := response.InternalServerError(err.Error())
				c.JSON(resp.Status, resp.Response)
				c.Abort()
				return
			}

			m.setUserInContext(c, user)
			c.Next()
			return
		} else if err != nil {
			resp := response.Unauthorized(err.Error())
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		if !user.IsVerify {
			resp := response.Unauthorized("Unverified User")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		if user.StatusAccount != entity.StatusAccountApproved {
			resp := response.Unauthorized("Account not approved")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		m.setUserInContext(c, user)

		c.Next()
	}
}

func (m *Middleware) setUserInContext(c *gin.Context, user *entity.User) {
	c.Set("user_id", user.ID)
	c.Set("user_role", user.Role)
	c.Set("user_email", user.Email)
}

func (m *Middleware) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			resp := response.Unauthorized("User role not found")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		if fmt.Sprintf("%v", role) != string(entity.UserRoleAdmin) {
			resp := response.Unauthorized("Admin access required")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *Middleware) StaffMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			resp := response.Unauthorized("User role not found")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		if fmt.Sprintf("%v", role) != string(entity.UserRoleStaff) {
			resp := response.Unauthorized("Staff access required")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *Middleware) StudentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			resp := response.Unauthorized("User role not found")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		if fmt.Sprintf("%v", role) != string(entity.UserRoleStudent) {
			resp := response.Unauthorized("Student access required")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *Middleware) ManagerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			resp := response.Unauthorized("User role not found")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		if fmt.Sprintf("%v", role) != string(entity.UserRoleAdmin) && fmt.Sprintf("%v", role) != string(entity.UserRoleStaff) {
			resp := response.Unauthorized("Manager access required")
			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *Middleware) ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			m.logger.Error("API error", zap.Error(err))

			var resp response.StatusResponse

			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				resp = m.formatValidationError(validationErrors)
				c.JSON(resp.Status, resp.Response)
				c.Abort()
				return
			}

			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				resp = response.NotFound("Resource not found")
			default:
				resp = response.InternalServerError("Internal server error")
			}

			c.JSON(resp.Status, resp.Response)
			c.Abort()
			return
		}
	}
}

func (m *Middleware) CustomRecovery() gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		resp := response.InternalServerError("Internal server error")
		c.JSON(resp.Status, resp.Response)
	}
}

func (m *Middleware) formatValidationError(validationErrors validator.ValidationErrors) response.StatusResponse {
	errors := make(map[string]string)

	for _, err := range validationErrors {
		field := err.Field()
		tag := err.Tag()

		var message string
		switch tag {
		case "required":
			message = field + " is required"
		case "email":
			message = field + " must be a valid email"
		case "min":
			message = field + " must be at least " + err.Param() + " characters"
		case "max":
			message = field + " must be at most " + err.Param() + " characters"
		case "len":
			message = field + " must be exactly " + err.Param() + " characters"
		case "numeric":
			message = field + " must be a number"
		case "alpha":
			message = field + " must contain only letters"
		case "alphanum":
			message = field + " must contain only letters and numbers"
		case "oneof":
			message = field + " must be one of the following values: " + err.Param()
		case "validdate":
			message = field + " must be a valid date"
		case "gtefield":
			message = field + " must be greater than or equal to " + err.Param()
		default:
			message = field + " is invalid"
		}

		fieldName := helper.ToSnakeCase(field)
		errors[fieldName] = message
	}

	return response.Validation("Invalid data", errors)
}
