package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type BillRepository interface {
	GetBill(ctx context.Context, bill *entity.Bill) error
	CreateBill(ctx context.Context) error
	MyListBills(ctx context.Context, userID uint64, query *entity.QueryBill) ([]entity.Bill, int64, error)
	UpdateBillsWithPayment(ctx context.Context, billIDs []uint64, paymentID uint64) error
	GetListBillsByIDs(ctx context.Context, billIDs []uint64) ([]entity.Bill, error)
}
