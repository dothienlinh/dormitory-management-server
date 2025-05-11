package services

import (
	"dormitory_management/internal/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db           *gorm.DB
	redis        *redis.Client
	logger       *zap.Logger
	Contract     *ContractService
	User         *UserService
	Room         *RoomService
	RoomCategory *RoomCategoryService
	Auth         *AuthService
}

func NewService(db *gorm.DB, redis *redis.Client, logger *zap.Logger, util *utils.Util) *Service {
	contractService := NewContractService(db, redis, logger)
	userService := NewUserService(db, redis, logger)
	roomService := NewRoomService(db, redis, logger)
	roomCategoryService := NewRoomCategoryService(db, redis, logger)
	authService := NewAuthService(db, redis, logger, util)

	return &Service{
		db:           db,
		redis:        redis,
		logger:       logger,
		Contract:     contractService,
		User:         userService,
		Room:         roomService,
		RoomCategory: roomCategoryService,
		Auth:         authService,
	}
}
