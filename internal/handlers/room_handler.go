package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ValidateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("id")

		room := models.Room{}
		if err := h.dbClient.Where("id = ?", roomId).First(&room).Error; err != nil {
			h.logger.Error("Error getting room", zap.Error(err))
			NotFound(c, "Room not found")
			return
		}

		c.Set("room", room)
		c.Next()
	}
}

func (h *Handler) CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := models.CreateRoom{}

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

		roomCategoryId := payload.RoomCategoryID
		roomCategory := models.RoomCategory{}
		if err := h.dbClient.Where("id = ?", roomCategoryId).First(&roomCategory).Error; err != nil {
			h.logger.Error("Error getting room category", zap.Error(err))
			NotFound(c, "Room category not found")
			return
		}

		if roomCategory.ID == 0 {
			h.logger.Error("Room category not found")
			NotFound(c, "Room category not found")
			return
		}

		roomNumber := payload.RoomNumber
		room := models.Room{}
		if err := h.dbClient.Where("room_number = ?", roomNumber).First(&room).Error; err != nil {
			h.logger.Error("Error getting room", zap.Error(err))
			BadRequest(c, errors.New("error getting room").Error())
			return
		}

		if room.ID != 0 {
			h.logger.Error("Room number already exists")
			BadRequest(c, "Room number already exists")
			return
		}

		roomCreate := models.Room{
			RoomNumber:     payload.RoomNumber,
			Status:         payload.Status,
			RoomCategoryID: roomCategory.ID,
		}
		if err := h.dbClient.Create(&roomCreate).Error; err != nil {
			h.logger.Error("Error creating room", zap.Error(err))
			BadRequest(c, errors.New("error creating room").Error())
			return
		}

		Success(c, CreateResponse{ID: roomCreate.ID}, 0)
	}
}

func (h *Handler) GetRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := models.FilterRoom{}

		if err := c.ShouldBindQuery(&query); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		query.Parse()

		rooms := []models.Room{}
		queryBuilder := h.dbClient.Model(&models.Room{})

		conditions := []string{}
		values := []interface{}{}

		if query.RoomNumber != "" {
			conditions = append(conditions, "room_number ILIKE ?")
			values = append(values, "%"+query.RoomNumber+"%")
		}

		if query.Status != "" {
			conditions = append(conditions, "status = ?")
			values = append(values, query.Status)
		}

		if query.RoomCategoryID != 0 {
			conditions = append(conditions, "room_category_id = ?")
			values = append(values, query.RoomCategoryID)
		}

		whereClause := strings.Join(conditions, " AND ")

		if err := queryBuilder.Select("id").Where(whereClause, values...).Count(&query.Total).Error; err != nil {
			h.logger.Error("Error getting rooms", zap.Error(err))
			BadRequest(c, errors.New("error getting rooms").Error())
			return
		}

		if err := queryBuilder.Select("*").Preload("RoomCategory").Limit(query.Limit).Offset(query.GetOffset()).Find(&rooms).Error; err != nil {
			h.logger.Error("Error getting rooms", zap.Error(err))
			BadRequest(c, errors.New("error getting rooms").Error())
			return
		}

		Success(c, rooms, query.Total)
	}
}

func (h *Handler) GetRoomDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("id")
		room := models.Room{}
		if err := h.dbClient.Where("id = ?", roomId).Preload("RoomCategory").First(&room).Error; err != nil {
			h.logger.Error("Error getting room", zap.Error(err))
			NotFound(c, "Room not found")
			return
		}

		Success(c, room, 0)
	}
}

func (h *Handler) UpdateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		room, exists := GetDataFromContext[models.Room](c, "room")
		if !exists {
			h.logger.Error("Room not found")
			NotFound(c, "Room not found")
			return
		}

		payload := models.UpdateRoom{}

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

		room.Status = payload.Status
		room.RoomCategoryID = payload.RoomCategoryID
		if err := h.dbClient.Save(&room).Error; err != nil {
			h.logger.Error("Error updating room", zap.Error(err))
			BadRequest(c, errors.New("error updating room").Error())
			return
		}

		Success(c, nil, 0)
	}
}

func (h *Handler) DeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		room, exists := GetDataFromContext[models.Room](c, "room")
		if !exists {
			h.logger.Error("Room not found")
			NotFound(c, "Room not found")
			return
		}

		if err := h.dbClient.Delete(&room).Error; err != nil {
			h.logger.Error("Error deleting room", zap.Error(err))
			BadRequest(c, errors.New("error deleting room").Error())
			return
		}

		Success(c, nil, 0)
	}
}
