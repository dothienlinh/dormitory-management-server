package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ValidateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("id")

		room := models.Room{}
		if err := h.dbClient.Where("id = ?", roomId).First(&room).Error; err != nil {
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
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		roomCategoryId := payload.RoomCategoryID
		roomCategory := models.RoomCategory{}
		h.dbClient.Where("id = ?", roomCategoryId).First(&roomCategory)

		if roomCategory.ID == 0 {
			NotFound(c, "Room category not found")
			return
		}

		roomNumber := payload.RoomNumber
		room := models.Room{}
		h.dbClient.Where("room_number = ?", roomNumber).First(&room)

		if room.ID != 0 {
			BadRequest(c, "Room number already exists")
			return
		}

		roomCreate := models.Room{
			RoomNumber:     payload.RoomNumber,
			Status:         payload.Status,
			RoomCategoryID: roomCategory.ID,
		}
		h.dbClient.Create(&roomCreate)

		Success(c, CreateResponse{ID: roomCreate.ID}, 0)
	}
}

func (h *Handler) GetRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := models.FilterRoom{}

		if err := c.ShouldBindQuery(&query); err != nil {
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

		queryString := strings.Join(conditions, " AND ")

		queryBuilder.Select("id").Where(queryString, values...).Count(&query.Total)

		queryBuilder.Select("*").Preload("RoomCategory").Limit(query.Limit).Offset(query.GetOffset()).Find(&rooms)

		Success(c, rooms, query.Total)
	}
}

func (h *Handler) GetRoomDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("id")
		room := models.Room{}
		if err := h.dbClient.Where("id = ?", roomId).Preload("RoomCategory").First(&room).Error; err != nil {
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
			NotFound(c, "Room not found")
			return
		}

		payload := models.UpdateRoom{}

		if err := c.ShouldBindJSON(&payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		room.Status = payload.Status
		room.RoomCategoryID = payload.RoomCategoryID
		h.dbClient.Save(&room)

		Success(c, nil, 0)
	}
}

func (h *Handler) DeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		room, exists := GetDataFromContext[models.Room](c, "room")
		if !exists {
			NotFound(c, "Room not found")
			return
		}

		h.dbClient.Delete(&room)

		Success(c, nil, 0)
	}
}
