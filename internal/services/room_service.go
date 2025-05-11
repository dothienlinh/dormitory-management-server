package services

import (
	"dormitory_management/internal/models"
	"errors"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoomService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

func NewRoomService(db *gorm.DB, redis *redis.Client, logger *zap.Logger) *RoomService {
	return &RoomService{db: db, redis: redis, logger: logger}
}

func (s *RoomService) ValidateRoom(roomID uint) (models.Room, error) {
	room := models.Room{}
	if err := s.db.Where("id = ?", roomID).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		return models.Room{}, err
	}

	return room, nil
}

func (s *RoomService) CreateRoom(payload models.CreateRoom) (*models.Room, error) {
	roomCategoryId := payload.RoomCategoryID
	roomCategory := models.RoomCategory{}
	if err := s.db.Where("id = ?", roomCategoryId).First(&roomCategory).Error; err != nil {
		s.logger.Error("Error getting room category", zap.Error(err))
		return nil, err
	}

	if roomCategory.ID == 0 {
		s.logger.Error("Room category not found")
		return nil, errors.New("room category not found")
	}

	roomNumber := payload.RoomNumber
	room := models.Room{}
	if err := s.db.Where("room_number = ?", roomNumber).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		return nil, err
	}

	if room.ID != 0 {
		s.logger.Error("Room number already exists")
		return nil, errors.New("room number already exists")
	}

	roomCreate := &models.Room{
		RoomNumber:     payload.RoomNumber,
		Status:         payload.Status,
		RoomCategoryID: roomCategory.ID,
	}
	if err := s.db.Create(roomCreate).Error; err != nil {
		s.logger.Error("Error creating room", zap.Error(err))
		return nil, err
	}

	return roomCreate, nil
}

func (s *RoomService) GetRooms(filter models.FilterRoom) ([]models.Room, error) {
	filter.Parse()

	rooms := []models.Room{}
	queryBuilder := s.db.Model(&models.Room{})

	whereClause, values := filter.Build()

	if err := queryBuilder.Select("id").Where(whereClause, values...).Count(&filter.Total).Error; err != nil {
		s.logger.Error("Error getting rooms", zap.Error(err))
		return nil, err
	}

	if err := queryBuilder.Select("*").Preload("RoomCategory").Limit(filter.Limit).Offset(filter.GetOffset()).Find(&rooms).Error; err != nil {
		s.logger.Error("Error getting rooms", zap.Error(err))
		return nil, err
	}

	return rooms, nil
}

func (s *RoomService) GetRoomDetail(roomID uint) (models.Room, error) {
	room := models.Room{}
	if err := s.db.Where("id = ?", roomID).Preload("RoomCategory").First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		return models.Room{}, err
	}

	return room, nil
}

func (s *RoomService) UpdateRoom(room models.Room, payload models.UpdateRoom) error {
	room.Status = payload.Status
	room.RoomCategoryID = payload.RoomCategoryID
	if err := s.db.Model(&room).Updates(room).Error; err != nil {
		s.logger.Error("Error updating room", zap.Error(err))
		return err
	}

	return nil
}

func (s *RoomService) DeleteRoom(room models.Room) error {
	if err := s.db.Delete(&room).Error; err != nil {
		s.logger.Error("Error deleting room", zap.Error(err))
		return err
	}

	return nil
}

func (s *RoomService) GetListStudentsInRoom(roomID uint) ([]models.UserSimple, error) {
	room := models.Room{}
	if err := s.db.Where("id = ?", roomID).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		return nil, err
	}

	users := []models.UserSimple{}
	if err := s.db.Model(&models.User{}).
		Select("id, full_name, student_code, email, gender, status, phone, birthday, avatar").
		Where("room_id = ?", roomID).Find(&users).Error; err != nil {
		s.logger.Error("Error getting users", zap.Error(err))
		return nil, err
	}

	return users, nil
}
