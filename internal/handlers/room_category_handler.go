package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ValidateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategoryId := c.Param("id")

		roomCategory := models.RoomCategory{}
		if err := h.dbClient.Where("id = ?", roomCategoryId).First(&roomCategory).Error; err != nil {
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

		roomCategory := models.RoomCategory{
			Name:        payload.Name,
			Capacity:    payload.Capacity,
			Price:       payload.Price,
			Acreage:     payload.Acreage,
			Description: payload.Description,
		}

		if err := h.dbClient.Create(&roomCategory).Error; err != nil {
			h.logger.Error("Error creating room category", zap.Error(err))
			BadRequest(c, errors.New("error creating room category").Error())
			return
		}

		Success(c, CreateResponse{ID: roomCategory.ID}, 0)
	}
}

func (h *Handler) GetRoomCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategories := []models.ListRoomCategory{}
		if err := h.dbClient.Model(&models.RoomCategory{}).Find(&roomCategories).Error; err != nil {
			h.logger.Error("Error getting room categories", zap.Error(err))
			BadRequest(c, errors.New("error getting room categories").Error())
			return
		}

		Success(c, roomCategories, int64(len(roomCategories)))
	}
}

func (h *Handler) GetRoomCategoryDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategoryId := c.Param("id")

		roomCategory := models.RoomCategory{}
		if err := h.dbClient.Where("id = ?", roomCategoryId).First(&roomCategory).Error; err != nil {
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

		roomCategory.Name = payload.Name
		roomCategory.Capacity = payload.Capacity
		roomCategory.Price = payload.Price
		roomCategory.Acreage = payload.Acreage
		roomCategory.Description = payload.Description

		if err := h.dbClient.Save(&roomCategory).Error; err != nil {
			h.logger.Error("Error updating room category", zap.Error(err))
			BadRequest(c, errors.New("error updating room category").Error())
			return
		}

		Success(c, nil, 0)
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

		if err := h.dbClient.Delete(&roomCategory).Error; err != nil {
			h.logger.Error("Error deleting room category", zap.Error(err))
			BadRequest(c, errors.New("error deleting room category").Error())
			return
		}

		Success(c, nil, 0)
	}
}
