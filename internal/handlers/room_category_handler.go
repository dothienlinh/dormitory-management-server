package handlers

import (
	"dormitory_management/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ValidateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategoryId := c.Param("id")
		roomCategoryIdUint, err := strconv.Atoi(roomCategoryId)
		if err != nil {
			h.logger.Error("Error converting roomCategoryId to uint", zap.Error(err))
			h.response.BadRequest(c, "Invalid roomCategoryId")
			return
		}

		roomCategory, err := h.service.RoomCategory.ValidateRoomCategory(uint(roomCategoryIdUint))
		if err != nil {
			h.logger.Error("Error getting room category", zap.Error(err))
			h.response.NotFound(c, "Room category not found")
			return
		}

		c.Set("roomCategory", roomCategory)
		c.Next()
	}
}

func (h *Handler) CreateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := models.CreateRoomCategory{}

		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.RoomCategory.CreateRoomCategory(c, payload)
	}
}

func (h *Handler) GetRoomCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := models.FilterRoomCategory{}
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		roomCategories, err := h.service.RoomCategory.GetRoomCategories(filter)
		if err != nil {
			h.logger.Error("Error getting room categories", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.response.Success(c, roomCategories, int64(len(roomCategories)))
	}
}

func (h *Handler) GetRoomCategoryDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategoryId := c.Param("id")
		roomCategoryIdUint, err := strconv.Atoi(roomCategoryId)
		if err != nil {
			h.logger.Error("Error converting roomCategoryId to uint", zap.Error(err))
			h.response.BadRequest(c, "Invalid roomCategoryId")
			return
		}

		h.service.RoomCategory.GetRoomCategoryDetail(c, uint(roomCategoryIdUint))
	}
}

func (h *Handler) UpdateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategory, exists := GetDataFromContext[models.RoomCategory](c, "roomCategory")
		if !exists {
			h.logger.Error("Room category not found")
			h.response.NotFound(c, "Room category not found")
			return
		}

		payload := models.UpdateRoomCategory{}

		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.RoomCategory.UpdateRoomCategory(c, roomCategory, payload)
	}
}

func (h *Handler) DeleteRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategory, exists := GetDataFromContext[models.RoomCategory](c, "roomCategory")
		if !exists {
			h.logger.Error("Room category not found")
			h.response.NotFound(c, "Room category not found")
			return
		}

		h.service.RoomCategory.DeleteRoomCategory(c, roomCategory)
	}
}
