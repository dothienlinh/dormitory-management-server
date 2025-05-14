package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RoomHandler handles HTTP requests related to rooms
type RoomHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

// NewRoomHandler creates a new RoomHandler
func NewRoomHandler(useCases usecase.UseCases, logger logger.Logger) *RoomHandler {
	return &RoomHandler{
		useCases: useCases,
		logger:   logger,
	}
}

// CreateRoom handles the request to create a room
func (h *RoomHandler) CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("CreateRoom")

		var room entity.Room
		if err := c.ShouldBindJSON(&room); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Room().CreateRoom(c, &room)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetRoomByID handles the request to get a room by ID
func (h *RoomHandler) GetRoomByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("GetRoomByID")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid room ID",
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Room().GetRoomByID(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}

// GetListRooms handles the request to get a list of rooms
func (h *RoomHandler) GetListRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("GetListRooms")

		var filter entity.RoomFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Room().GetListRooms(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

// UpdateRoom handles the request to update a room
func (h *RoomHandler) UpdateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("UpdateRoom")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid room ID",
				"error":   "Bad Request",
			})
			return
		}

		var room entity.Room
		if err := c.ShouldBindJSON(&room); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Room().UpdateRoom(c, uint(id), &room)
		c.JSON(resp.Status, resp.Response)
	}
}

// DeleteRoom handles the request to delete a room
func (h *RoomHandler) DeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("DeleteRoom")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid room ID",
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Room().DeleteRoom(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}

// RoomCategoryHandler handles HTTP requests related to room categories
type RoomCategoryHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

// NewRoomCategoryHandler creates a new RoomCategoryHandler
func NewRoomCategoryHandler(useCases usecase.UseCases, logger logger.Logger) *RoomCategoryHandler {
	return &RoomCategoryHandler{
		useCases: useCases,
		logger:   logger,
	}
}

// CreateRoomCategory handles the request to create a room category
func (h *RoomCategoryHandler) CreateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("CreateRoomCategory")

		var category entity.RoomCategory
		if err := c.ShouldBindJSON(&category); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.RoomCategory().CreateRoomCategory(c, &category)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetRoomCategoryByID handles the request to get a room category by ID
func (h *RoomCategoryHandler) GetRoomCategoryByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("GetRoomCategoryByID")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid room category ID",
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.RoomCategory().GetRoomCategoryByID(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}

// GetListRoomCategories handles the request to get a list of room categories
func (h *RoomCategoryHandler) GetListRoomCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("GetListRoomCategories")

		var filter entity.RoomCategoryFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.RoomCategory().GetListRoomCategories(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

// UpdateRoomCategory handles the request to update a room category
func (h *RoomCategoryHandler) UpdateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("UpdateRoomCategory")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid room category ID",
				"error":   "Bad Request",
			})
			return
		}

		var category entity.RoomCategory
		if err := c.ShouldBindJSON(&category); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.RoomCategory().UpdateRoomCategory(c, uint(id), &category)
		c.JSON(resp.Status, resp.Response)
	}
}

// DeleteRoomCategory handles the request to delete a room category
func (h *RoomCategoryHandler) DeleteRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("DeleteRoomCategory")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid room category ID",
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.RoomCategory().DeleteRoomCategory(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}
