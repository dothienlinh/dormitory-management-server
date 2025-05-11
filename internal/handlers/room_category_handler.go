package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
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
			BadRequest(c, "Invalid roomCategoryId")
			return
		}

		roomCategory, err := h.service.RoomCategory.ValidateRoomCategory(uint(roomCategoryIdUint))
		if err != nil {
			h.logger.Error("Error getting room category", zap.Error(err))
			NotFound(c, "Room category not found")
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
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			h.logger.Error("Error validating struct", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		roomCategory, err := h.service.RoomCategory.CreateRoomCategory(payload)
		if err != nil {
			h.logger.Error("Error creating room category", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, CreateResponse{ID: roomCategory.ID}, 0)
	}
}

func (h *Handler) GetRoomCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := models.FilterRoomCategory{}
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		roomCategories, err := h.service.RoomCategory.GetRoomCategories(filter)
		if err != nil {
			h.logger.Error("Error getting room categories", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, roomCategories, int64(len(roomCategories)))
	}
}

func (h *Handler) GetRoomCategoryDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategoryId := c.Param("id")
		roomCategoryIdUint, err := strconv.Atoi(roomCategoryId)
		if err != nil {
			h.logger.Error("Error converting roomCategoryId to uint", zap.Error(err))
			BadRequest(c, "Invalid roomCategoryId")
			return
		}

		roomCategory, err := h.service.RoomCategory.GetRoomCategoryDetail(uint(roomCategoryIdUint))
		if err != nil {
			h.logger.Error("Error getting room category", zap.Error(err))
			NotFound(c, "Room category not found")
			return
		}

		Success(c, roomCategory, 0)
	}
}

func (h *Handler) UpdateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategory, exists := GetDataFromContext[models.RoomCategory](c, "roomCategory")
		if !exists {
			h.logger.Error("Room category not found")
			NotFound(c, "Room category not found")
			return
		}

		payload := models.UpdateRoomCategory{}

		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			h.logger.Error("Error validating struct", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		roomCategory, err := h.service.RoomCategory.UpdateRoomCategory(roomCategory, payload)
		if err != nil {
			h.logger.Error("Error updating room category", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, roomCategory, 0)
	}
}

func (h *Handler) DeleteRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategory, exists := GetDataFromContext[models.RoomCategory](c, "roomCategory")
		if !exists {
			h.logger.Error("Room category not found")
			NotFound(c, "Room category not found")
			return
		}

		err := h.service.RoomCategory.DeleteRoomCategory(roomCategory)
		if err != nil {
			h.logger.Error("Error deleting room category", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, nil, 0)
	}
}
