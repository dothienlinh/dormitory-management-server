package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ValidateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategoryId := c.Param("id")

		roomCategory := models.RoomCategory{}
		if err := h.dbClient.Where("id = ?", roomCategoryId).First(&roomCategory).Error; err != nil {
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
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
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

		h.dbClient.Create(&roomCategory)

		Success(c, CreateResponse{ID: roomCategory.ID}, 0)
	}
}

func (h *Handler) GetRoomCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategories := []models.ListRoomCategory{}
		h.dbClient.Model(&models.RoomCategory{}).Find(&roomCategories)
		Success(c, roomCategories, int64(len(roomCategories)))
	}
}

func (h *Handler) GetRoomCategoryDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategoryId := c.Param("id")

		roomCategoryTemp := models.RoomCategory{}
		if err := h.dbClient.Where("id = ?", roomCategoryId).First(&roomCategoryTemp).Error; err != nil {
			NotFound(c, "Room category not found")
			return
		}

		roomCategory := models.RoomCategoryDetail{
			BaseModel:   roomCategoryTemp.BaseModel,
			Name:        roomCategoryTemp.Name,
			Capacity:    roomCategoryTemp.Capacity,
			Price:       roomCategoryTemp.Price,
			Acreage:     roomCategoryTemp.Acreage,
			Description: roomCategoryTemp.Description,
		}

		var simpleRooms []models.RoomSimple
		h.dbClient.Model(&models.Room{}).Select("id, created_at, updated_at, deleted_at, room_number, status").
			Where("room_category_id = ?", roomCategoryId).
			Find(&simpleRooms)

		roomCategory.Rooms = simpleRooms

		Success(c, roomCategory, 0)
	}
}

func (h *Handler) UpdateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategory, exists := GetDataFromContext[models.RoomCategory](c, "roomCategory")
		if !exists {
			NotFound(c, "Room category not found")
			return
		}

		payload := models.UpdateRoomCategory{}

		if err := c.ShouldBindJSON(&payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		roomCategory.Name = payload.Name
		roomCategory.Capacity = payload.Capacity
		roomCategory.Price = payload.Price
		roomCategory.Acreage = payload.Acreage
		roomCategory.Description = payload.Description

		h.dbClient.Save(&roomCategory)
		Success(c, nil, 0)
	}
}

func (h *Handler) DeleteRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomCategory, exists := GetDataFromContext[models.RoomCategory](c, "roomCategory")
		if !exists {
			NotFound(c, "Room category not found")
			return
		}

		h.dbClient.Delete(&roomCategory)
		Success(c, nil, 0)
	}
}
