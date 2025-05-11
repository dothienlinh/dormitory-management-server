package services

import (
	"dormitory_management/internal/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	BaseService struct {
		db       *gorm.DB
		redis    *redis.Client
		logger   *zap.Logger
		response *APIResponse
	}

	Service struct {
		Contract     *ContractService
		User         *UserService
		Room         *RoomService
		RoomCategory *RoomCategoryService
		Auth         *AuthService
	}

	APIResponse struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Data       any    `json:"data,omitempty"`
		Errors     any    `json:"errors,omitempty"`
		Total      int64  `json:"total,omitempty"`
	}
)

func NewService(db *gorm.DB, redis *redis.Client, logger *zap.Logger, util *utils.Util, response *APIResponse) *Service {
	contractService := NewContractService(db, redis, logger, response)
	userService := NewUserService(db, redis, logger, response)
	roomService := NewRoomService(db, redis, logger, response)
	roomCategoryService := NewRoomCategoryService(db, redis, logger, response)
	authService := NewAuthService(db, redis, logger, util, response)

	return &Service{
		Contract:     contractService,
		User:         userService,
		Room:         roomService,
		RoomCategory: roomCategoryService,
		Auth:         authService,
	}
}

func NewAPIResponse() *APIResponse {
	return &APIResponse{}
}

func (response *APIResponse) Success(c *gin.Context, data any, total int64) {
	c.JSON(http.StatusOK, APIResponse{
		StatusCode: http.StatusOK,
		Data:       data,
		Total:      total,
		Message:    "success",
	})
}

func (response *APIResponse) InternalServer(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	})
}

func (response *APIResponse) Unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	})
}

func (response *APIResponse) BadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	})
}

func (response *APIResponse) Forbidden(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
		StatusCode: http.StatusForbidden,
		Message:    message,
	})
}

func (response *APIResponse) NotFound(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, APIResponse{
		StatusCode: http.StatusNotFound,
		Message:    message,
	})
}

func (response *APIResponse) DBError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.NotFound(c, err.Error())
	} else {
		response.InternalServer(c, err.Error())
	}
}
