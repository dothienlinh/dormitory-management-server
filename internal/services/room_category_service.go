package services

import (
	"dormitory_management/internal/models"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoomCategoryService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

func NewRoomCategoryService(db *gorm.DB, redis *redis.Client, logger *zap.Logger) *RoomCategoryService {
	return &RoomCategoryService{db: db, redis: redis, logger: logger}
}

func (s *RoomCategoryService) ValidateRoomCategory(roomCategoryID uint) (models.RoomCategory, error) {
	roomCategory := models.RoomCategory{}
	if err := s.db.Where("id = ?", roomCategoryID).First(&roomCategory).Error; err != nil {
		s.logger.Error("Error getting room category", zap.Error(err))
		return models.RoomCategory{}, err
	}

	return roomCategory, nil
}

func (s *RoomCategoryService) GetRoomCategories(filter models.FilterRoomCategory) ([]models.ListRoomCategory, error) {
	filter.Parse()
	whereClause, values := filter.Build()

	roomCategories := []models.ListRoomCategory{}
	if err := s.db.Model(&models.RoomCategory{}).Where(whereClause, values...).Find(&roomCategories).Error; err != nil {
		s.logger.Error("Error getting room categories", zap.Error(err))
		return nil, err
	}

	return roomCategories, nil
}

func (s *RoomCategoryService) CreateRoomCategory(roomCategory models.CreateRoomCategory) (models.RoomCategory, error) {
	roomCategoryCreate := models.RoomCategory{
		Name:        roomCategory.Name,
		Capacity:    roomCategory.Capacity,
		Price:       roomCategory.Price,
		Acreage:     roomCategory.Acreage,
		Description: roomCategory.Description,
	}

	if err := s.db.Create(&roomCategoryCreate).Error; err != nil {
		s.logger.Error("Error creating room category", zap.Error(err))
		return models.RoomCategory{}, err
	}

	return roomCategoryCreate, nil
}

func (s *RoomCategoryService) GetRoomCategoryDetail(roomCategoryID uint) (models.RoomCategory, error) {
	roomCategory := models.RoomCategory{}

	if err := s.db.Where("id = ?", roomCategoryID).First(&roomCategory).Error; err != nil {
		s.logger.Error("Error getting room category", zap.Error(err))
		return models.RoomCategory{}, err
	}

	return roomCategory, nil
}

func (s *RoomCategoryService) UpdateRoomCategory(roomCategory models.RoomCategory, payload models.UpdateRoomCategory) (models.RoomCategory, error) {
	roomCategory.Name = payload.Name
	roomCategory.Capacity = payload.Capacity
	roomCategory.Price = payload.Price
	roomCategory.Acreage = payload.Acreage
	roomCategory.Description = payload.Description

	if err := s.db.Model(&roomCategory).Updates(roomCategory).Error; err != nil {
		s.logger.Error("Error updating room category", zap.Error(err))
		return models.RoomCategory{}, err
	}

	return roomCategory, nil
}

func (s *RoomCategoryService) DeleteRoomCategory(roomCategory models.RoomCategory) error {
	if err := s.db.Delete(&roomCategory).Error; err != nil {
		s.logger.Error("Error deleting room category", zap.Error(err))
		return err
	}

	return nil
}
