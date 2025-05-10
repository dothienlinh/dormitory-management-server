package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
	"strconv"
	"strings"

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

		filter.Parse()

		conditions := []string{"role = ?"}
		values := []interface{}{models.UserRoleStudent}

		if filter.Status != "" {
			conditions = append(conditions, "status = ?")
			values = append(values, filter.Status)
		}

		if filter.Gender != "" {
			conditions = append(conditions, "gender = ?")
			values = append(values, filter.Gender)
		}

		if filter.Keyword != "" {
			conditions = append(conditions, "(full_name LIKE ? OR email LIKE ? OR phone LIKE ? OR student_code LIKE ?)")
			keyword := "%" + filter.Keyword + "%"
			values = append(values, keyword, keyword, keyword, keyword)
		}

		whereClause := strings.Join(conditions, " AND ")
		h.logger.Info("whereClause", zap.String("whereClause", whereClause))
		h.logger.Info("values", zap.Any("values", values))
		users := []models.User{}

		query := h.dbClient.Model(&models.User{}).Preload("Room")

		if err := query.Where(whereClause, values...).Count(&filter.Total).Error; err != nil {
			h.logger.Error("Error querying data", zap.Error(err))
			BadRequest(c, "Error querying data")
			return
		}

		if err := query.Where(whereClause, values...).
			Order("created_at DESC").
			Limit(filter.Limit).
			Offset(filter.GetOffset()).
			Find(&users).Error; err != nil {
			h.logger.Error("Error querying data", zap.Error(err))
			BadRequest(c, "Error querying data")
			return
		}

		h.logger.Info("GetListUser success")

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

		room := models.Room{}
		if err := h.dbClient.Preload("RoomCategory").Where("id = ?", roomID).First(&room).Error; err != nil {
			h.logger.Error("Error getting room", zap.Error(err))
			BadRequest(c, "Room not found")
			return
		}

		var countStudentInRoom int64
		if err := h.dbClient.Model(&models.RoomRent{}).Select("id").Where("room_id = ?", room.ID).Count(&countStudentInRoom).Error; err != nil {
			h.logger.Error("Error getting count student in room", zap.Error(err))
			DBError(c, err)
			return
		}

		if countStudentInRoom >= int64(room.RoomCategory.Capacity) {
			BadRequest(c, "Room is full")
			return
		}

		user := models.User{}
		if err := h.dbClient.Where("id = ? AND role = ?", userID, models.UserRoleStudent).First(&user).Error; err != nil {
			h.logger.Error("Error getting user", zap.Error(err))
			BadRequest(c, "Student not found")
			return
		}

		if user.RoomRentID != nil {
			BadRequest(c, "Student already has a room")
			return
		}

		roomRent := models.RoomRent{
			RoomID: uint(roomID),
			UserID: uint(userID),
			Status: payload.Status,
		}

		if err := h.dbClient.Create(&roomRent).Error; err != nil {
			h.logger.Error("Error creating room rent", zap.Error(err))
			DBError(c, err)
			return
		}

		if err := h.dbClient.Model(&user).Update("room_rent_id", roomRent.ID).Error; err != nil {
			h.logger.Error("Error updating user", zap.Error(err))
			DBError(c, err)
			return
		}

		h.logger.Info("======> Student added to room successfully")

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

		room := models.Room{}
		if err := h.dbClient.Where("id = ?", roomID).First(&room).Error; err != nil {
			h.logger.Error("Error getting room", zap.Error(err))
			BadRequest(ctx, "Room not found")
			return
		}

		user := models.User{}
		if err := h.dbClient.Preload("RoomRent").Where("id = ? AND role = ?", userID, models.UserRoleStudent).First(&user).Error; err != nil {
			h.logger.Error("Error getting user", zap.Error(err))
			BadRequest(ctx, "Student not found")
			return
		}

		if user.RoomRentID == nil {
			BadRequest(ctx, "Student not in room")
			return
		}

		if err := h.dbClient.Model(&user).Update("room_rent_id", nil).Error; err != nil {
			h.logger.Error("Error removing user from room", zap.Error(err))
			DBError(ctx, err)
			return
		}

		if err := h.dbClient.Delete(&user.RoomRent).Error; err != nil {
			h.logger.Error("Error deleting room rent", zap.Error(err))
			DBError(ctx, err)
			return
		}

		h.logger.Info("======> Student removed from room successfully")

		Success(ctx, "Student removed from room successfully", 0)
	}
}
