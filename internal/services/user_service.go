package services

import (
	"dormitory_management/internal/models"
	"errors"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

func NewUserService(db *gorm.DB, redis *redis.Client, logger *zap.Logger) *UserService {
	return &UserService{db: db, redis: redis, logger: logger}
}

func (s *UserService) GetListUser(filter *models.FilterUser) ([]models.User, error) {
	users := []models.User{}

	filter.Parse()
	whereClause, values := filter.Build()

	s.logger.Info("whereClause", zap.String("whereClause", whereClause))
	s.logger.Info("values", zap.Any("values", values))

	query := s.db.Model(&models.User{}).Preload("Room")

	if err := query.Where(whereClause, values...).Count(&filter.Total).Error; err != nil {
		s.logger.Error("Error querying data", zap.Error(err))
		return nil, err
	}

	if err := query.Where(whereClause, values...).
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&users).Error; err != nil {
		s.logger.Error("Error querying data", zap.Error(err))
		return nil, err
	}

	s.logger.Info("GetListUser success", zap.Any("users", users))

	return users, nil
}

func (s *UserService) AddUserToRoom(userID uint, roomID uint, payload models.CreateRoomRent) error {
	room := models.Room{}
	if err := s.db.Preload("RoomCategory").Where("id = ?", roomID).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		return err
	}

	var countStudentInRoom int64
	if err := s.db.Model(&models.RoomRent{}).Select("id").Where("room_id = ?", room.ID).Count(&countStudentInRoom).Error; err != nil {
		s.logger.Error("Error getting count student in room", zap.Error(err))
		return err
	}

	if countStudentInRoom >= int64(room.RoomCategory.Capacity) {
		return errors.New("room is full")
	}

	user := models.User{}
	if err := s.db.Where("id = ? AND role = ?", userID, models.UserRoleStudent).First(&user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		return err
	}

	if user.RoomRentID != nil {
		return errors.New("student already has a room")
	}

	roomRent := models.RoomRent{
		RoomID: uint(roomID),
		UserID: uint(userID),
		Status: payload.Status,
	}

	if err := s.db.Create(&roomRent).Error; err != nil {
		s.logger.Error("Error creating room rent", zap.Error(err))
		return err
	}

	if err := s.db.Model(&user).Update("room_rent_id", roomRent.ID).Error; err != nil {
		s.logger.Error("Error updating user", zap.Error(err))
		return err
	}

	s.logger.Info("======> Student added to room successfully")

	return nil
}

func (s *UserService) RemoveUserFromRoom(userID uint, roomID uint) error {
	room := models.Room{}
	if err := s.db.Where("id = ?", roomID).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		return err
	}

	user := models.User{}
	if err := s.db.Preload("RoomRent").Where("id = ? AND role = ?", userID, models.UserRoleStudent).First(&user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		return err
	}

	if user.RoomRentID == nil {
		return errors.New("student not in room")
	}

	if err := s.db.Model(&user).Update("room_rent_id", nil).Error; err != nil {
		s.logger.Error("Error removing user from room", zap.Error(err))
		return err
	}

	if err := s.db.Delete(&user.RoomRent).Error; err != nil {
		s.logger.Error("Error deleting room rent", zap.Error(err))
		return err
	}

	return nil
}
