package usecase

import (
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"

	"github.com/hibiken/asynq"
)

// useCases implements the usecase.UseCases interface
type useCases struct {
	user         usecase.UserUseCase
	room         usecase.RoomUseCase
	roomCategory usecase.RoomCategoryUseCase
	contract     usecase.ContractUseCase
	auth         usecase.AuthUseCase
	dashboard    usecase.DashboardUseCase
	email        usecase.EmailUseCase
}

// NewUseCases creates a new UseCases instance
func NewUseCases(repos repository.Repositories, logger logger.Logger, asynqClient *asynq.Client) usecase.UseCases {
	useCases := &useCases{}

	useCases.user = NewUserUseCase(repos, logger)
	useCases.room = NewRoomUseCase(repos, logger)
	useCases.roomCategory = NewRoomCategoryUseCase(repos, logger)
	useCases.contract = NewContractUseCase(repos, logger)
	useCases.auth = NewAuthUseCase(repos, logger)
	useCases.dashboard = NewDashboardUseCase(repos, logger)
	useCases.email = NewEmailUseCase(repos, logger, asynqClient)

	return useCases
}

// User returns the user use case
func (uc *useCases) User() usecase.UserUseCase {
	return uc.user
}

// Room returns the room use case
func (uc *useCases) Room() usecase.RoomUseCase {
	return uc.room
}

// RoomCategory returns the room category use case
func (uc *useCases) RoomCategory() usecase.RoomCategoryUseCase {
	return uc.roomCategory
}

// Contract returns the contract use case
func (uc *useCases) Contract() usecase.ContractUseCase {
	return uc.contract
}

// Auth returns the auth use case
func (uc *useCases) Auth() usecase.AuthUseCase {
	return uc.auth
}

// Dashboard returns the dashboard use case
func (uc *useCases) Dashboard() usecase.DashboardUseCase {
	return uc.dashboard
}

func (uc *useCases) Email() usecase.EmailUseCase {
	return uc.email
}
