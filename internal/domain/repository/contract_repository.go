package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

// ContractRepository defines the interface for contract data operations
type ContractRepository interface {
	// Create a new contract
	Create(ctx context.Context, contract *entity.Contract) error

	// GetByID retrieves a contract by ID
	GetByID(ctx context.Context, id uint) (*entity.Contract, error)

	// GetByUserID retrieves a contract by user ID
	GetByUserID(ctx context.Context, userID uint) (*entity.Contract, error)

	// List retrieves contracts based on filter
	List(ctx context.Context, filter *entity.ContractFilter) ([]entity.Contract, int64, error)

	// Update updates an existing contract
	Update(ctx context.Context, contract *entity.Contract) error

	// Delete deletes a contract by ID
	Delete(ctx context.Context, id uint) error
}
