package handlers

import (
	"dormitory_management/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ValidateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("id")
		roomIdUint, err := strconv.Atoi(roomId)
		if err != nil {
			h.logger.Error("Error converting roomID to uint", zap.Error(err))
			h.response.BadRequest(c, "Invalid roomID")
			return
		}

		h.service.Room.ValidateRoom(c, uint(roomIdUint))

	}
}

func (h *Handler) CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := models.CreateRoom{}

		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.Room.CreateRoom(c, payload)
	}
}

func (h *Handler) GetRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := models.FilterRoom{}

		if err := c.ShouldBindQuery(&query); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.Room.GetRooms(c, query)
	}
}

func (h *Handler) GetRoomDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("id")
		roomIdUint, err := strconv.Atoi(roomId)
		if err != nil {
			h.logger.Error("Error converting roomID to uint", zap.Error(err))
			h.response.BadRequest(c, "Invalid roomID")
			return
		}

		h.service.Room.GetRoomDetail(c, uint(roomIdUint))
	}
}

func (h *Handler) UpdateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		room, exists := GetDataFromContext[models.Room](c, "room")
		if !exists {
			h.logger.Error("Room not found")
			h.response.NotFound(c, "Room not found")
			return
		}

		payload := models.UpdateRoom{}

		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.Room.UpdateRoom(c, room, payload)

	}
}

func (h *Handler) DeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		room, exists := GetDataFromContext[models.Room](c, "room")
		if !exists {
			h.logger.Error("Room not found")
			h.response.NotFound(c, "Room not found")
			return
		}

		h.service.Room.DeleteRoom(c, room)

	}
}

func (h *Handler) GetListStudentsInRoom() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomIdParam := ctx.Param("id")
		roomId, err := strconv.Atoi(roomIdParam)
		if err != nil {
			h.logger.Error("Error converting roomID to int", zap.Error(err))
			h.response.BadRequest(ctx, "Invalid roomID")
			return
		}

		h.service.Room.GetListStudentsInRoom(ctx, uint(roomId))

	}
}
