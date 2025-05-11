package services

import (
	"dormitory_management/internal/models"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ContractService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

func NewContractService(db *gorm.DB, redis *redis.Client, logger *zap.Logger) *ContractService {
	return &ContractService{db: db, redis: redis, logger: logger}
}

func (s *ContractService) CreateContract(payload models.CreateContract) (*models.Contract, error) {
	contract := &models.Contract{
		RoomID:    payload.RoomID,
		StartDate: payload.StartDate,
		EndDate:   payload.EndDate,
		Status:    payload.Status,
		UserId:    payload.UserId,
	}

	if err := s.db.Create(contract).Error; err != nil {
		s.logger.Error("failed to create contract", zap.Error(err))
		return nil, err
	}

	return contract, nil
}
