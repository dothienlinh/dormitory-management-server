package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type contractTermRepository struct {
	db *gorm.DB
}

func NewContractTermRepository(db *gorm.DB) repository.ContractTermRepository {
	return &contractTermRepository{
		db: db,
	}
}

func (r *contractTermRepository) Create(ctx context.Context, contractTerm *entity.ContractTerm) error {
	if err := r.db.WithContext(ctx).Table(entity.ContractTerm{}.TableName()).Create(contractTerm).Error; err != nil {
		return fmt.Errorf("failed to create contract term: %w", err)
	}
	return nil
}

func (r *contractTermRepository) GetByID(ctx context.Context, id uint) (*entity.ContractTerm, error) {
	var contractTerm entity.ContractTerm
	if err := r.db.WithContext(ctx).Table(contractTerm.TableName()).
		Preload("Contract").
		First(&contractTerm, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contract term with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get contract term: %w", err)
	}
	return &contractTerm, nil
}

func (r *contractTermRepository) GetByContractID(ctx context.Context, contractID uint) ([]entity.ContractTerm, error) {
	var contractTerms []entity.ContractTerm
	if err := r.db.WithContext(ctx).Table(entity.ContractTerm{}.TableName()).
		Where("contract_id = ?", contractID).
		Preload("Contract").
		Order("created_at ASC").
		Find(&contractTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to get contract terms for contract %d: %w", contractID, err)
	}
	return contractTerms, nil
}

func (r *contractTermRepository) List(ctx context.Context, contractID uint) ([]entity.ContractTerm, error) {
	var contractTerms []entity.ContractTerm
	query := r.db.WithContext(ctx).Table(entity.ContractTerm{}.TableName())

	if contractID > 0 {
		query = query.Where("contract_id = ?", contractID)
	}

	if err := query.Preload("Contract").
		Order("created_at ASC").
		Find(&contractTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to get contract terms: %w", err)
	}

	return contractTerms, nil
}

func (r *contractTermRepository) Update(ctx context.Context, contractTerm *entity.ContractTerm) error {
	if err := r.db.WithContext(ctx).Table(contractTerm.TableName()).Save(contractTerm).Error; err != nil {
		return fmt.Errorf("failed to update contract term: %w", err)
	}
	return nil
}

func (r *contractTermRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Table(entity.ContractTerm{}.TableName()).Delete(&entity.ContractTerm{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete contract term: %w", err)
	}
	return nil
}
