package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type paymentHistoryRepository struct {
	db *gorm.DB
}

func NewPaymentHistoryRepository(db *gorm.DB) repository.PaymentHistoryRepository {
	return &paymentHistoryRepository{
		db: db,
	}
}

func (r *paymentHistoryRepository) Create(ctx context.Context, paymentHistory *entity.PaymentHistory) error {
	if err := r.db.WithContext(ctx).Table(entity.PaymentHistory{}.TableName()).Create(paymentHistory).Error; err != nil {
		return fmt.Errorf("failed to create payment history: %w", err)
	}
	return nil
}

func (r *paymentHistoryRepository) GetByID(ctx context.Context, id uint) (*entity.PaymentHistory, error) {
	var paymentHistory entity.PaymentHistory
	if err := r.db.WithContext(ctx).Table(paymentHistory.TableName()).
		Preload("Contract").
		Preload("Contract.User").
		Preload("Contract.Room.RoomCategory").
		First(&paymentHistory, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment history with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get payment history: %w", err)
	}
	return &paymentHistory, nil
}

func (r *paymentHistoryRepository) GetByContractID(ctx context.Context, contractID uint) ([]entity.PaymentHistory, error) {
	var paymentHistories []entity.PaymentHistory
	if err := r.db.WithContext(ctx).Table(entity.PaymentHistory{}.TableName()).
		Where("contract_id = ?", contractID).
		Preload("Contract").
		Preload("Contract.User").
		Preload("Contract.Room.RoomCategory").
		Order("due_date DESC").
		Find(&paymentHistories).Error; err != nil {
		return nil, fmt.Errorf("failed to get payment histories for contract %d: %w", contractID, err)
	}
	return paymentHistories, nil
}

func (r *paymentHistoryRepository) List(ctx context.Context, filter *entity.PaymentHistoryFilter) ([]entity.PaymentHistory, int64, error) {
	filter.Parse()

	query := r.db.WithContext(ctx).Table(entity.PaymentHistory{}.TableName())

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.ContractID != 0 {
		query = query.Where("contract_id = ?", filter.ContractID)
	}

	if filter.PaymentType != "" {
		query = query.Where("payment_type = ?", filter.PaymentType)
	}

	if filter.Keyword != "" {
		query = query.Joins("JOIN contracts ON payment_histories.contract_id = contracts.id").
			Joins("JOIN users ON contracts.user_id = users.id").
			Where("users.full_name LIKE ? OR users.email LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count payment histories: %w", err)
	}

	var paymentHistories []entity.PaymentHistory
	if err := query.Preload("Contract").
		Preload("Contract.User").
		Preload("Contract.Room.RoomCategory").
		Order("due_date DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&paymentHistories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get payment histories: %w", err)
	}

	return paymentHistories, total, nil
}

func (r *paymentHistoryRepository) GetMyPaymentHistory(ctx context.Context, userID uint64, filter *entity.PaymentHistoryFilter) ([]entity.PaymentHistory, int64, error) {
	filter.Parse()

	query := r.db.WithContext(ctx).Table(entity.PaymentHistory{}.TableName()).
		Joins("JOIN contracts ON payment_histories.contract_id = contracts.id").
		Where("contracts.user_id = ?", userID)

	if filter.Status != "" {
		query = query.Where("payment_histories.status = ?", filter.Status)
	}

	if filter.PaymentType != "" {
		query = query.Where("payment_histories.payment_type = ?", filter.PaymentType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count payment histories for user %d: %w", userID, err)
	}

	var paymentHistories []entity.PaymentHistory
	if err := query.Preload("Contract").
		Preload("Contract.User").
		Preload("Contract.Room.RoomCategory").
		Order("payment_histories.due_date DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&paymentHistories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get payment histories for user %d: %w", userID, err)
	}

	return paymentHistories, total, nil
}

func (r *paymentHistoryRepository) Update(ctx context.Context, paymentHistory *entity.PaymentHistory) error {
	if err := r.db.WithContext(ctx).Table(paymentHistory.TableName()).Save(paymentHistory).Error; err != nil {
		return fmt.Errorf("failed to update payment history: %w", err)
	}
	return nil
}

func (r *paymentHistoryRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Table(entity.PaymentHistory{}.TableName()).Delete(&entity.PaymentHistory{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete payment history: %w", err)
	}
	return nil
}
