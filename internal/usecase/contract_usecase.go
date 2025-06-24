package usecase

import (
	"context"
	"dormitory_management/internal/common"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contractUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

func NewContractUseCase(repos repository.Repositories, logger logger.Logger) usecase.ContractUseCase {
	return &contractUseCase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *contractUseCase) CreateContract(ctx context.Context, createContract *entity.CreateContract) response.StatusResponse {
	user := &entity.User{
		Base: entity.Base{ID: createContract.UserID},
	}
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to validate user", zap.Error(err))
		return response.BadRequest("Invalid user")
	}

	room, err := uc.repos.Room().GetByID(ctx, createContract.RoomID)
	if err != nil {
		uc.logger.Error("Failed to validate room", zap.Error(err))
		return response.BadRequest("Invalid room")
	}
	var amountStudents int64
	if err := uc.repos.Room().AmountStudentsInRoom(ctx, room.ID, &amountStudents); err != nil {
		uc.logger.Error("Failed to check room capacity", zap.Error(err))
		return response.BadRequest("Room capacity exceeded")
	}

	if amountStudents >= int64(room.RoomCategory.Capacity) {
		uc.logger.Error("Room capacity exceeded", zap.Int64("current", amountStudents), zap.Int("capacity", room.RoomCategory.Capacity))
		return response.BadRequest("Room capacity exceeded")
	}

	startDate, err := common.ParsedTime(createContract.StartDate)
	if err != nil {
		uc.logger.Error("Invalid start date format", zap.Error(err))
		return response.InternalServerError("Invalid start date format")
	}

	endDate, err := common.ParsedTime(createContract.EndDate)
	if err != nil {
		uc.logger.Error("Invalid start date format", zap.Error(err))
		return response.InternalServerError("Invalid start date format")
	}

	contract := entity.Contract{
		UserID:      createContract.UserID,
		RoomID:      createContract.RoomID,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: createContract.Description,
		Status:      entity.ContractStatusInactive,
		Price:       createContract.Price,
		Code:        uuid.New().String(),
	}

	if createContract.Price == 0 {
		contract.Price = room.RoomCategory.Price
	}

	if time.Now().After(startDate) {
		contract.Status = entity.ContractStatusActive
	}

	if err := uc.repos.Contract().Create(ctx, &contract); err != nil {
		uc.logger.Error("Failed to create contract", zap.Error(err))
		return response.InternalServerError("Failed to create contract")
	}

	return response.Created(createContract)
}

func (uc *contractUseCase) GetContractByID(ctx context.Context, id uint) response.StatusResponse {
	contract, err := uc.repos.Contract().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get contract by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Contract with ID %d not found", id))
	}

	return response.Success(contract, 1)
}

func (uc *contractUseCase) GetContractByUserID(ctx context.Context, userID uint) response.StatusResponse {
	contract, err := uc.repos.Contract().GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get contract by user ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Contract for user ID %d not found", userID))
	}

	return response.Success(contract, 1)
}

func (uc *contractUseCase) GetListContracts(ctx context.Context, filter *entity.ContractFilter) response.StatusResponse {
	contracts, total, err := uc.repos.Contract().List(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to get list of contracts", zap.Error(err))
		return response.InternalServerError("Failed to get list of contracts")
	}

	return response.Success(contracts, total)
}

func (uc *contractUseCase) UpdateContract(ctx context.Context, id uint, updateData *entity.UpdateContract) response.StatusResponse {
	contract, err := uc.repos.Contract().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get contract for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Contract with ID %d not found", id))
	}

	if updateData.Status != "" {
		contract.Status = updateData.Status
	}
	if updateData.StartDate != nil {
		contract.StartDate = *updateData.StartDate
	}
	if updateData.EndDate != nil {
		contract.EndDate = *updateData.EndDate
	}
	if updateData.Price != nil {
		contract.Price = *updateData.Price
	}
	if updateData.Description != nil {
		contract.Description = *updateData.Description
	}

	if err := uc.repos.Contract().Update(ctx, contract); err != nil {
		uc.logger.Error("Failed to update contract", zap.Error(err))
		return response.InternalServerError("Failed to update contract")
	}

	return response.Success(contract, 1)
}

func (uc *contractUseCase) DeleteContract(ctx context.Context, id uint) response.StatusResponse {
	_, err := uc.repos.Contract().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get contract for deletion", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Contract with ID %d not found", id))
	}

	if err := uc.repos.Contract().Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete contract", zap.Error(err))
		return response.InternalServerError("Failed to delete contract")
	}

	return response.Success("Contract deleted successfully", 0)
}
