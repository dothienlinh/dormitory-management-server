package handlers

import (
	"dormitory_management/internal/services"
	"dormitory_management/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	Handler struct {
		response *services.APIResponse
		service  *services.Service
		logger   *zap.Logger
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

func NewHandler(service *services.Service, logger *zap.Logger, response *services.APIResponse) *Handler {

	return &Handler{service: service, logger: logger, response: response}
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
