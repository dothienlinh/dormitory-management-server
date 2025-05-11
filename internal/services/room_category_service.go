package services

import (
	"dormitory_management/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoomCategoryService struct {
	BaseService
}

func NewRoomCategoryService(db *gorm.DB, redis *redis.Client, logger *zap.Logger, response *APIResponse) *RoomCategoryService {
	return &RoomCategoryService{BaseService{db: db, redis: redis, logger: logger, response: response}}
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

func (s *RoomCategoryService) CreateRoomCategory(c *gin.Context, roomCategory models.CreateRoomCategory) {
	roomCategoryCreate := models.RoomCategory{
		Name:        roomCategory.Name,
		Capacity:    roomCategory.Capacity,
		Price:       roomCategory.Price,
		Acreage:     roomCategory.Acreage,
		Description: roomCategory.Description,
	}

	if err := s.db.Create(&roomCategoryCreate).Error; err != nil {
		s.logger.Error("Error creating room category", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, roomCategoryCreate, 0)
}

func (s *RoomCategoryService) GetRoomCategoryDetail(c *gin.Context, roomCategoryID uint) {
	roomCategory := models.RoomCategory{}

	if err := s.db.Where("id = ?", roomCategoryID).First(&roomCategory).Error; err != nil {
		s.logger.Error("Error getting room category", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, roomCategory, 0)
}

func (s *RoomCategoryService) UpdateRoomCategory(c *gin.Context, roomCategory models.RoomCategory, payload models.UpdateRoomCategory) {
	roomCategory.Name = payload.Name
	roomCategory.Capacity = payload.Capacity
	roomCategory.Price = payload.Price
	roomCategory.Acreage = payload.Acreage
	roomCategory.Description = payload.Description

	if err := s.db.Model(&roomCategory).Updates(roomCategory).Error; err != nil {
		s.logger.Error("Error updating room category", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, roomCategory, 0)
}

func (s *RoomCategoryService) DeleteRoomCategory(c *gin.Context, roomCategory models.RoomCategory) {
	if err := s.db.Delete(&roomCategory).Error; err != nil {
		s.logger.Error("Error deleting room category", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, "Room category deleted successfully", 0)
}
