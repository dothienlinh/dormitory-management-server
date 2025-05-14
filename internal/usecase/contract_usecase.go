package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

// contractUseCase implements the usecase.ContractUseCase interface
type contractUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

// NewContractUseCase creates a new contract use case
func NewContractUseCase(repos repository.Repositories, logger logger.Logger) usecase.ContractUseCase {
	return &contractUseCase{
		repos:  repos,
		logger: logger,
	}
}

// CreateContract creates a new contract
func (uc *contractUseCase) CreateContract(ctx context.Context, contract *entity.Contract) response.StatusResponse {
	// Validate user
	user, err := uc.repos.User().GetByID(ctx, contract.UserID)
	if err != nil {
		uc.logger.Error("Failed to validate user", zap.Error(err))
		return response.BadRequest("Invalid user")
	}

	// Validate room
	room, err := uc.repos.Room().GetByID(ctx, contract.RoomID)
	if err != nil {
		uc.logger.Error("Failed to validate room", zap.Error(err))
		return response.BadRequest("Invalid room")
	}

	// Check if user already has a contract
	if user.ContractID != nil {
		return response.BadRequest("User already has an active contract")
	}

	// Set contract price from room category if not specified
	if contract.Price == 0 {
		contract.Price = room.RoomCategory.Price
	}

	if err := uc.repos.Contract().Create(ctx, contract); err != nil {
		uc.logger.Error("Failed to create contract", zap.Error(err))
		return response.InternalServerError("Failed to create contract")
	}

	// Update user with contract ID
	user.ContractID = &contract.ID
	if err := uc.repos.User().Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user with contract", zap.Error(err))
		return response.InternalServerError("Failed to update user with contract")
	}

	return response.Created(contract)
}

// GetContractByID retrieves a contract by ID
func (uc *contractUseCase) GetContractByID(ctx context.Context, id uint) response.StatusResponse {
	contract, err := uc.repos.Contract().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get contract by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Contract with ID %d not found", id))
	}

	return response.Success(contract, 1)
}

// GetContractByUserID retrieves a contract by user ID
func (uc *contractUseCase) GetContractByUserID(ctx context.Context, userID uint) response.StatusResponse {
	contract, err := uc.repos.Contract().GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get contract by user ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Contract for user ID %d not found", userID))
	}

	return response.Success(contract, 1)
}

// GetListContracts retrieves contracts based on filter
func (uc *contractUseCase) GetListContracts(ctx context.Context, filter *entity.ContractFilter) response.StatusResponse {
	contracts, total, err := uc.repos.Contract().List(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to get list of contracts", zap.Error(err))
		return response.InternalServerError("Failed to get list of contracts")
	}

	return response.Success(contracts, total)
}

// UpdateContract updates a contract
func (uc *contractUseCase) UpdateContract(ctx context.Context, id uint, updateData *entity.UpdateContract) response.StatusResponse {
	contract, err := uc.repos.Contract().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get contract for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Contract with ID %d not found", id))
	}

	// Update contract fields
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

// DeleteContract deletes a contract
func (uc *contractUseCase) DeleteContract(ctx context.Context, id uint) response.StatusResponse {
	contract, err := uc.repos.Contract().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get contract for deletion", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Contract with ID %d not found", id))
	}

	// Get user to remove contract reference
	user, err := uc.repos.User().GetByID(ctx, contract.UserID)
	if err == nil && user.ContractID != nil && *user.ContractID == id {
		user.ContractID = nil
		if err := uc.repos.User().Update(ctx, user); err != nil {
			uc.logger.Error("Failed to update user contract reference", zap.Error(err))
			// Continue with deletion anyway
		}
	}

	if err := uc.repos.Contract().Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete contract", zap.Error(err))
		return response.InternalServerError("Failed to delete contract")
	}

	return response.Success("Contract deleted successfully", 0)
}
