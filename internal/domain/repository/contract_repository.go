package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type ContractRepository interface {
	Create(ctx context.Context, contract *entity.Contract) error

	GetByID(ctx context.Context, id uint) (*entity.Contract, error)

	GetByUserID(ctx context.Context, userID uint) (*entity.Contract, error)

	List(ctx context.Context, filter *entity.ContractFilter) ([]entity.Contract, int64, error)

	Update(ctx context.Context, contract *entity.Contract) error

	Delete(ctx context.Context, id uint) error

	// New methods for student contract module
	GetMyContract(ctx context.Context, userID uint64) (*entity.Contract, error)

	GetContractWithRelations(ctx context.Context, id uint) (*entity.Contract, error)
}
