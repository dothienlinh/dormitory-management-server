package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/pkg/logger"
)

type BillUseCase struct {
	logger logger.Logger
	repo   repository.Repositories
}

func NewBillUseCase(logger logger.Logger, repo repository.Repositories) *BillUseCase {
	return &BillUseCase{
		logger: logger,
		repo:   repo,
	}
}

func (uc *BillUseCase) MyListBills(ctx context.Context, userID uint64, query *entity.QueryBill) ([]entity.Bill, int64, error) {
	return uc.repo.Bill().MyListBills(ctx, userID, query)
}
