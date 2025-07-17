package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type BillRepository interface {
	GetBill(ctx context.Context, bill *entity.Bill) error
}
