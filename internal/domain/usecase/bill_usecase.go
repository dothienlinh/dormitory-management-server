package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type BillUseCase interface {
	MyListBills(ctx context.Context, userID uint64, query *entity.QueryBill) ([]entity.Bill, int64, error)
}
