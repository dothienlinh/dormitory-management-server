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

type RoomHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewRoomHandler(useCases usecase.UseCases, logger logger.Logger) *RoomHandler {
	return &RoomHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *RoomHandler) CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("CreateRoom")

		var room entity.CreateRoom
		if err := c.ShouldBindJSON(&room); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Room().CreateRoom(c, &room)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomHandler) GetRoomByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetRoomByID")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			resp = response.BadRequest("Invalid room ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Room().GetRoomByID(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomHandler) GetListRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetListRooms")

		var filter entity.RoomFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Room().GetListRooms(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomHandler) UpdateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("UpdateRoom")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			resp = response.BadRequest("Invalid room ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var room entity.UpdateRoom
		if err := c.ShouldBindJSON(&room); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Room().UpdateRoom(c, id, &room)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomHandler) DeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("DeleteRoom")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			resp = response.BadRequest("Invalid room ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Room().DeleteRoom(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

type RoomCategoryHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewRoomCategoryHandler(useCases usecase.UseCases, logger logger.Logger) *RoomCategoryHandler {
	return &RoomCategoryHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *RoomCategoryHandler) CreateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("CreateRoomCategory")

		var category entity.CreateRoomCategory
		if err := c.ShouldBindJSON(&category); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.RoomCategory().CreateRoomCategory(c, &category)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomCategoryHandler) GetRoomCategoryByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetRoomCategoryByID")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.RoomCategory().GetRoomCategoryByID(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomCategoryHandler) GetListRoomCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetListRoomCategories")

		var filter entity.RoomCategoryFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.RoomCategory().GetListRoomCategories(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomCategoryHandler) UpdateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("UpdateRoomCategory")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var category entity.UpdateRoomCategory
		if err := c.ShouldBindJSON(&category); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.RoomCategory().UpdateRoomCategory(c, id, &category)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomCategoryHandler) DeleteRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("DeleteRoomCategory")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.RoomCategory().DeleteRoomCategory(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}
