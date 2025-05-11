package services

import (
	"dormitory_management/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	BaseService
}

func NewUserService(db *gorm.DB, redis *redis.Client, logger *zap.Logger, response *APIResponse) *UserService {
	return &UserService{BaseService{db: db, redis: redis, logger: logger, response: response}}
}

func (s *UserService) GetListUser(c *gin.Context, filter *models.FilterUser) {
	users := []models.User{}

	filter.Parse()
	whereClause, values := filter.Build()

	s.logger.Info("whereClause", zap.String("whereClause", whereClause))
	s.logger.Info("values", zap.Any("values", values))

	query := s.db.Model(&models.User{}).Preload("Room")

	if err := query.Where(whereClause, values...).Count(&filter.Total).Error; err != nil {
		s.logger.Error("Error querying data", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	if err := query.Where(whereClause, values...).
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&users).Error; err != nil {
		s.logger.Error("Error querying data", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.logger.Info("GetListUser success", zap.Any("users", users))

	s.response.Success(c, users, filter.Total)
}

func (s *UserService) AddUserToRoom(c *gin.Context, userID uint, roomID uint, payload models.CreateRoomRent) {
	room := models.Room{}
	if err := s.db.Preload("RoomCategory").Where("id = ?", roomID).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	var countStudentInRoom int64
	if err := s.db.Model(&models.RoomRent{}).Select("id").Where("room_id = ?", room.ID).Count(&countStudentInRoom).Error; err != nil {
		s.logger.Error("Error getting count student in room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	if countStudentInRoom >= int64(room.RoomCategory.Capacity) {
		s.response.BadRequest(c, "Room is full")
		return
	}

	user := models.User{}
	if err := s.db.Where("id = ? AND role = ?", userID, models.UserRoleStudent).First(&user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	if user.RoomRentID != nil {
		s.response.BadRequest(c, "Student already has a room")
		return
	}

	roomRent := models.RoomRent{
		RoomID: uint(roomID),
		UserID: uint(userID),
		Status: payload.Status,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&roomRent).Error; err != nil {
			s.logger.Error("Error creating room rent", zap.Error(err))
			return err
		}

		if err := tx.Model(&user).Update("room_rent_id", roomRent.ID).Error; err != nil {
			s.logger.Error("Error updating user", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		s.response.DBError(c, err)
		return
	}

	s.logger.Info("======> Student added to room successfully")

	s.response.Success(c, "Student added to room successfully", 0)
}

func (s *UserService) RemoveUserFromRoom(c *gin.Context, userID uint, roomID uint) {
	room := models.Room{}
	if err := s.db.Where("id = ?", roomID).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	user := models.User{}
	if err := s.db.Preload("RoomRent").Where("id = ? AND role = ?", userID, models.UserRoleStudent).First(&user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	if user.RoomRentID == nil {
		s.response.BadRequest(c, "Student not in room")
		return
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Update("room_rent_id", nil).Error; err != nil {
			s.logger.Error("Error removing user from room", zap.Error(err))
			return err
		}

		if err := tx.Delete(&user.RoomRent).Error; err != nil {
			s.logger.Error("Error deleting room rent", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, "Student removed from room successfully", 0)
}
