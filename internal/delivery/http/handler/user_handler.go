package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewUserHandler(useCases usecase.UseCases, logger logger.Logger) *UserHandler {
	return &UserHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		h.logger.Info("GetUserByID")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse user ID", zap.Error(err))
			resp = response.BadRequest("Invalid user ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.User().GetUserByID(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *UserHandler) GetListUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetListUsers")

		var filter entity.UserFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.User().GetListUsers(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("UpdateUser")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse user ID", zap.Error(err))
			resp = response.BadRequest("Invalid user ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var user entity.User
		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.User().UpdateUser(c, id, &user)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("DeleteUser")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse user ID", zap.Error(err))
			resp = response.BadRequest("Invalid user ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.User().DeleteUser(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *UserHandler) AddUserToRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("AddUserToRoom")

		var payload entity.AddUserToRoom
		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.User().AddUserToRoom(c, payload)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *UserHandler) UserLeavesRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("UserLeavesRoom")

		var payload entity.UserLeavesRoom

		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.User().UserLeavesRoom(c, payload)
		c.JSON(resp.Status, resp.Response)
	}
}
