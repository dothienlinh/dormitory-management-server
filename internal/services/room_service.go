package services

import (
	"dormitory_management/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoomService struct {
	BaseService
}

func NewRoomService(db *gorm.DB, redis *redis.Client, logger *zap.Logger, response *APIResponse) *RoomService {
	return &RoomService{BaseService{db: db, redis: redis, logger: logger, response: response}}
}

func (s *RoomService) ValidateRoom(c *gin.Context, roomID uint) {
	room := models.Room{}
	if err := s.db.Where("id = ?", roomID).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	c.Set("room", room)
	c.Next()
}

func (s *RoomService) CreateRoom(c *gin.Context, payload models.CreateRoom) {
	roomCategoryId := payload.RoomCategoryID
	roomCategory := models.RoomCategory{}
	if err := s.db.Where("id = ?", roomCategoryId).First(&roomCategory).Error; err != nil {
		s.logger.Error("Error getting room category", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	if roomCategory.ID == 0 {
		s.logger.Error("Room category not found")
		s.response.NotFound(c, "Room category not found")
		return
	}

	roomNumber := payload.RoomNumber
	room := models.Room{}
	if err := s.db.Where("room_number = ?", roomNumber).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	if room.ID != 0 {
		s.logger.Error("Room number already exists")
		s.response.BadRequest(c, "Room number already exists")
		return
	}

	roomCreate := &models.Room{
		RoomNumber:     payload.RoomNumber,
		Status:         payload.Status,
		RoomCategoryID: roomCategory.ID,
	}
	if err := s.db.Create(roomCreate).Error; err != nil {
		s.logger.Error("Error creating room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, roomCreate, 0)
}

func (s *RoomService) GetRooms(c *gin.Context, filter models.FilterRoom) {
	filter.Parse()

	rooms := []models.Room{}
	queryBuilder := s.db.Model(&models.Room{})

	whereClause, values := filter.Build()

	if err := queryBuilder.Select("id").Where(whereClause, values...).Count(&filter.Total).Error; err != nil {
		s.logger.Error("Error getting rooms", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	if err := queryBuilder.Select("*").Preload("RoomCategory").Limit(filter.Limit).Offset(filter.GetOffset()).Find(&rooms).Error; err != nil {
		s.logger.Error("Error getting rooms", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, rooms, filter.Total)
}

func (s *RoomService) GetRoomDetail(c *gin.Context, roomID uint) {
	room := models.Room{}
	if err := s.db.Where("id = ?", roomID).Preload("RoomCategory").First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, room, 0)
}

func (s *RoomService) UpdateRoom(c *gin.Context, room models.Room, payload models.UpdateRoom) {
	room.Status = payload.Status
	room.RoomCategoryID = payload.RoomCategoryID
	if err := s.db.Model(&room).Updates(room).Error; err != nil {
		s.logger.Error("Error updating room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, room, 0)
}

func (s *RoomService) DeleteRoom(c *gin.Context, room models.Room) {
	if err := s.db.Delete(&room).Error; err != nil {
		s.logger.Error("Error deleting room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, "Room deleted successfully", 0)
}

func (s *RoomService) GetListStudentsInRoom(c *gin.Context, roomID uint) {
	room := models.Room{}
	if err := s.db.Where("id = ?", roomID).First(&room).Error; err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	users := []models.UserSimple{}
	if err := s.db.Model(&models.User{}).
		Select("id, full_name, student_code, email, gender, status, phone, birthday, avatar").
		Where("room_id = ?", roomID).Find(&users).Error; err != nil {
		s.logger.Error("Error getting users", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, users, 0)
}
