package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"

	"gorm.io/gorm"
)

type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) repository.BillRepository {
	return &billRepository{db: db}
}

func (r *billRepository) GetBill(ctx context.Context, bill *entity.Bill) error {
	return r.db.WithContext(ctx).Table(bill.TableName()).Preload("User").Preload("Payment").First(bill).Error
}
