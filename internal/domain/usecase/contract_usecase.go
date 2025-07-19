package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type ContractUseCase interface {
	CreateContract(ctx context.Context, contract *entity.CreateContract) response.StatusResponse

	GetContractByID(ctx context.Context, id uint) response.StatusResponse

	GetContractByUserID(ctx context.Context, userID uint) response.StatusResponse

	GetListContracts(ctx context.Context, filter *entity.ContractFilter) response.StatusResponse

	UpdateContract(ctx context.Context, id uint, contract *entity.UpdateContract) response.StatusResponse

	DeleteContract(ctx context.Context, id uint) response.StatusResponse

	// New methods for student contract module
	GetMyContract(ctx context.Context, userID uint64) response.StatusResponse

	DownloadContractPDF(ctx context.Context, contractID uint, userID uint64) response.StatusResponse

	GetContractPaymentHistory(ctx context.Context, contractID uint, userID uint64) response.StatusResponse
}
