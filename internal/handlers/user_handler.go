package handlers

import (
	"dormitory_management/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetListUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("GetListUser")
		filter := models.FilterUser{}
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Error binding query", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.User.GetListUser(c, &filter)
	}
}

func (h *Handler) AddUserToRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("Add user to room")
		userIDParam := c.Param("user_id")
		roomIDParam := c.Param("room_id")

		userID, err := strconv.Atoi(userIDParam)
		if err != nil {
			h.logger.Error("Error converting userID to int", zap.Error(err))
			h.response.BadRequest(c, "Invalid userID")
			return
		}

		roomID, err := strconv.Atoi(roomIDParam)
		if err != nil {
			h.logger.Error("Error converting roomID to int", zap.Error(err))
			h.response.BadRequest(c, "Invalid roomID")
			return
		}

		payload := models.CreateRoomRent{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.User.AddUserToRoom(c, uint(userID), uint(roomID), payload)
	}
}

func (h *Handler) RemoveUserFromRoom() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.logger.Info("Remove user from room")
		userIDParam := ctx.Param("user_id")
		roomIDParam := ctx.Param("room_id")

		userID, err := strconv.Atoi(userIDParam)
		if err != nil {
			h.logger.Error("Error converting userID to int", zap.Error(err))
			h.response.BadRequest(ctx, "Invalid userID")
			return
		}

		roomID, err := strconv.Atoi(roomIDParam)
		if err != nil {
			h.logger.Error("Error converting roomID to int", zap.Error(err))
			h.response.BadRequest(ctx, "Invalid roomID")
			return
		}

		h.service.User.RemoveUserFromRoom(ctx, uint(userID), uint(roomID))
	}
}
