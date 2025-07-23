package tasks

import (
	"context"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/pkg/logger"
)

type BillTask struct {
	logger logger.Logger
	repos  repository.Repositories
}

func NewBillTask(logger logger.Logger, repos repository.Repositories) *BillTask {
	return &BillTask{
		logger: logger,
		repos:  repos,
	}
}

func (bt *BillTask) CreateBill() error {
	ctx := context.Background()
	return bt.repos.Bill().CreateBill(ctx)
}
