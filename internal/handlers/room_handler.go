package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
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
			BadRequest(c, "Invalid roomID")
			return
		}

		room, err := h.service.Room.ValidateRoom(uint(roomIdUint))
		if err != nil {
			h.logger.Error("Error validating room", zap.Error(err))
			BadRequest(c, err.Error())
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

		room, err := h.service.Room.CreateRoom(payload)
		if err != nil {
			h.logger.Error("Error creating room", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, CreateResponse{ID: room.ID}, 0)
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

		rooms, err := h.service.Room.GetRooms(query)
		if err != nil {
			h.logger.Error("Error getting rooms", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, rooms, query.Total)
	}
}

func (h *Handler) GetRoomDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("id")
		roomIdUint, err := strconv.Atoi(roomId)
		if err != nil {
			h.logger.Error("Error converting roomID to uint", zap.Error(err))
			BadRequest(c, "Invalid roomID")
			return
		}

		room, err := h.service.Room.GetRoomDetail(uint(roomIdUint))
		if err != nil {
			h.logger.Error("Error getting room", zap.Error(err))
			BadRequest(c, err.Error())
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

		if err := h.service.Room.UpdateRoom(room, payload); err != nil {
			h.logger.Error("Error updating room", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, "Updated room successfully", 0)
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

		if err := h.service.Room.DeleteRoom(room); err != nil {
			h.logger.Error("Error deleting room", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, "Deleted room successfully", 0)
	}
}

func (h *Handler) GetListStudentsInRoom() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomIdParam := ctx.Param("id")
		roomId, err := strconv.Atoi(roomIdParam)
		if err != nil {
			h.logger.Error("Error converting roomID to int", zap.Error(err))
			BadRequest(ctx, "Invalid roomID")
			return
		}

		users, err := h.service.Room.GetListStudentsInRoom(uint(roomId))
		if err != nil {
			h.logger.Error("Error getting list students in room", zap.Error(err))
			BadRequest(ctx, err.Error())
			return
		}

		Success(ctx, users, 0)
	}
}
