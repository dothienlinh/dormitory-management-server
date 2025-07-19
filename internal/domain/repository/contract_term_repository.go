package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type ContractTermRepository interface {
	Create(ctx context.Context, contractTerm *entity.ContractTerm) error

	GetByID(ctx context.Context, id uint) (*entity.ContractTerm, error)

	GetByContractID(ctx context.Context, contractID uint) ([]entity.ContractTerm, error)

	List(ctx context.Context, contractID uint) ([]entity.ContractTerm, error)

	Update(ctx context.Context, contractTerm *entity.ContractTerm) error

	Delete(ctx context.Context, id uint) error
}
