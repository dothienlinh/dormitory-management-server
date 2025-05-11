package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
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
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(filter); err != nil {
			h.logger.Error("Error validating filter", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		users, err := h.service.User.GetListUser(&filter)
		if err != nil {
			h.logger.Error("Error getting list user", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, users, filter.Total)
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
			BadRequest(c, "Invalid userID")
			return
		}

		roomID, err := strconv.Atoi(roomIDParam)
		if err != nil {
			h.logger.Error("Error converting roomID to int", zap.Error(err))
			BadRequest(c, "Invalid roomID")
			return
		}

		payload := models.CreateRoomRent{}
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

		if err := h.service.User.AddUserToRoom(uint(userID), uint(roomID), payload); err != nil {
			h.logger.Error("Error adding user to room", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, "Student added to room successfully", 0)
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
			BadRequest(ctx, "Invalid userID")
			return
		}

		roomID, err := strconv.Atoi(roomIDParam)
		if err != nil {
			h.logger.Error("Error converting roomID to int", zap.Error(err))
			BadRequest(ctx, "Invalid roomID")
			return
		}

		if err := h.service.User.RemoveUserFromRoom(uint(userID), uint(roomID)); err != nil {
			h.logger.Error("Error removing user from room", zap.Error(err))
			BadRequest(ctx, err.Error())
			return
		}

		Success(ctx, "Student removed from room successfully", 0)
	}
}
