package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"time"

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

func (r *billRepository) CreateBill(ctx context.Context) error {
	contracts := []entity.Contract{}
	if err := r.db.WithContext(ctx).Table(entity.Contract{}.TableName()).
		Preload("Room.RoomCategory").
		Where(
			"NOW() - INTERVAL '1 MONTH' BETWEEN start_date AND end_date AND status = ?",
			entity.ContractStatusActive,
		).
		Find(&contracts).Error; err != nil {
		return err
	}

	bills := []entity.Bill{}
	for _, v := range contracts {
		bills = append(bills, entity.Bill{
			UserID:  v.UserID,
			Amount:  v.Room.RoomCategory.Price,
			Status:  entity.BillStatusPending,
			DueDate: time.Now().Add(15 * 24 * time.Hour), // Due date is set to 15 days from now
		})
	}

	if len(bills) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Table(entity.Bill{}.TableName()).Create(&bills).Error; err != nil {
		return err
	}

	return nil
}

func (r *billRepository) MyListBills(ctx context.Context, userID uint64, query *entity.QueryBill) ([]entity.Bill, int64, error) {
	bills := []entity.Bill{}
	queryBuilder := r.db.WithContext(ctx).Table(entity.Bill{}.TableName()).Where("user_id = ?", userID)

	if query.Description != nil {
		queryBuilder = queryBuilder.Where("description ILIKE ?", "%"+*query.Description+"%")
	}

	if query.Status != nil {
		queryBuilder = queryBuilder.Where("status = ?", *query.Status)
	}

	total := int64(0)
	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, total, err
	}

	if err := queryBuilder.
		Limit(query.Limit).
		Offset(query.GetOffset()).
		Find(&bills).Error; err != nil {
		return nil, total, err
	}

	return bills, total, nil
}

func (r *billRepository) UpdateBillsWithPayment(ctx context.Context, billIDs []uint64, paymentID uint64) error {
	return r.db.WithContext(ctx).Table(entity.Bill{}.TableName()).Where("id IN ?", billIDs).Updates(entity.Bill{PaymentID: &paymentID}).Error
}

func (r *billRepository) GetListBillsByIDs(ctx context.Context, billIDs []uint64) ([]entity.Bill, error) {
	bills := []entity.Bill{}
	if err := r.db.WithContext(ctx).Table(entity.Bill{}.TableName()).Where("id IN ? AND (status = ? OR payment_id IS NOT NULL)", billIDs, entity.BillStatusPaid).Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}
