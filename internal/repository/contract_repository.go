package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type contractRepository struct {
	db *gorm.DB
}

func NewContractRepository(db *gorm.DB) repository.ContractRepository {
	return &contractRepository{
		db: db,
	}
}

func (r *contractRepository) Create(ctx context.Context, createContract *entity.CreateContract) error {
	if err := r.db.WithContext(ctx).Table(entity.Contract{}.TableName()).Create(createContract).Error; err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}
	return nil
}

func (r *contractRepository) GetByID(ctx context.Context, id uint) (*entity.Contract, error) {
	var contract entity.Contract
	if err := r.db.WithContext(ctx).Table(contract.TableName()).Preload("User").Preload("Room.RoomCategory").First(&contract, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contract with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}
	return &contract, nil
}

func (r *contractRepository) GetByUserID(ctx context.Context, userID uint) (*entity.Contract, error) {
	var contract entity.Contract
	if err := r.db.WithContext(ctx).Table(contract.TableName()).Where("user_id = ?", userID).Preload("User").Preload("Room.RoomCategory").First(&contract).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contract for user ID %d not found", userID)
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}
	return &contract, nil
}

func (r *contractRepository) List(ctx context.Context, filter *entity.ContractFilter) ([]entity.Contract, int64, error) {
	filter.Parse()

	query := r.db.WithContext(ctx).Table(entity.Contract{}.TableName())

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.RoomID > 0 {
		query = query.Where("room_id = ?", filter.RoomID)
	}

	if filter.Keyword != "" {
		query = query.Where("description LIKE ?", "%"+filter.Keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	var contracts []entity.Contract
	if err := query.Preload("User").Preload("Room.RoomCategory").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&contracts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get contracts: %w", err)
	}

	return contracts, total, nil
}

func (r *contractRepository) Update(ctx context.Context, contract *entity.Contract) error {
	if err := r.db.WithContext(ctx).Table(contract.TableName()).Save(contract).Error; err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}
	return nil
}

func (r *contractRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Table(entity.Contract{}.TableName()).Delete(&entity.Contract{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}
	return nil
}
