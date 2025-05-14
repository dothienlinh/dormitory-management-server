package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

// ContractUseCase defines the interface for contract business logic
type ContractUseCase interface {
	// CreateContract creates a new contract
	CreateContract(ctx context.Context, contract *entity.Contract) response.StatusResponse

	// GetContractByID retrieves a contract by ID
	GetContractByID(ctx context.Context, id uint) response.StatusResponse

	// GetContractByUserID retrieves a contract by user ID
	GetContractByUserID(ctx context.Context, userID uint) response.StatusResponse

	// GetListContracts retrieves contracts based on filter
	GetListContracts(ctx context.Context, filter *entity.ContractFilter) response.StatusResponse

	// UpdateContract updates a contract
	UpdateContract(ctx context.Context, id uint, contract *entity.UpdateContract) response.StatusResponse

	// DeleteContract deletes a contract
	DeleteContract(ctx context.Context, id uint) response.StatusResponse
}
