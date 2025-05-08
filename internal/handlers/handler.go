package handlers

import (
	"dormitory_management/internal/utils"
	"dormitory_management/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type (
	Handler struct {
		dbClient    *gorm.DB
		redisClient *redis.Client
		util        *utils.Util
	}

	APIResponse struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Data       any    `json:"data,omitempty"`
		Errors     any    `json:"errors,omitempty"`
		Total      int64  `json:"total,omitempty"`
	}

	CreateResponse struct {
		ID uint `json:"id"`
	}
)

func NewHandler(dbClient *gorm.DB, redisClient *redis.Client, util *utils.Util) *Handler {
	return &Handler{dbClient: dbClient, redisClient: redisClient, util: util}
}

func GetDataFromContext[T any](c *gin.Context, key string) (value T, exists bool) {
	data, exists := c.Get(key)
	if !exists {
		return *new(T), false
	}

	value, ok := data.(T)
	if !ok {
		return *new(T), false
	}

	return value, true
}

func Success(c *gin.Context, data any, total int64) {
	c.JSON(http.StatusOK, APIResponse{
		StatusCode: http.StatusOK,
		Data:       data,
		Total:      total,
		Message:    "success",
	})
}

func InternalServer(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		StatusCode: http.StatusForbidden,
		Message:    message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		StatusCode: http.StatusNotFound,
		Message:    message,
	})
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, e := range c.Errors {
			switch err := e.Err.(type) {
			case *pkg.HttpError:
				res := map[string]any{"message": err.Message}
				if err.Errors != nil {
					res["errors"] = err.Errors
				}
				c.JSON(err.StatusCode, res)
				return
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
				return
			}
		}
	}
}
