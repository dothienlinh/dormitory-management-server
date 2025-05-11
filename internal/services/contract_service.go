package services

import (
	"dormitory_management/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ContractService struct {
	BaseService
}

func NewContractService(db *gorm.DB, redis *redis.Client, logger *zap.Logger, response *APIResponse) *ContractService {
	return &ContractService{BaseService{db: db, redis: redis, logger: logger, response: response}}
}

func (s *ContractService) CreateContract(ctx *gin.Context, payload models.CreateContract) {
	contract := &models.Contract{
		RoomID:    payload.RoomID,
		StartDate: payload.StartDate,
		EndDate:   payload.EndDate,
		Status:    payload.Status,
		UserId:    payload.UserId,
	}

	if err := s.db.Create(contract).Error; err != nil {
		s.logger.Error("failed to create contract", zap.Error(err))
		s.response.DBError(ctx, err)
		return
	}

	s.response.Success(ctx, contract, 0)
}
