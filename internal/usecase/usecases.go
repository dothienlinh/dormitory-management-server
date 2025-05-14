package usecase

import (
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
)

// useCases implements the usecase.UseCases interface
type useCases struct {
	repos        repository.Repositories
	logger       logger.Logger
	user         usecase.UserUseCase
	room         usecase.RoomUseCase
	roomCategory usecase.RoomCategoryUseCase
	contract     usecase.ContractUseCase
	auth         usecase.AuthUseCase
}

// NewUseCases creates a new UseCases instance
func NewUseCases(repos repository.Repositories, logger logger.Logger) usecase.UseCases {
	useCases := &useCases{
		repos:  repos,
		logger: logger,
	}

	useCases.user = NewUserUseCase(repos, logger)
	useCases.room = NewRoomUseCase(repos, logger)
	useCases.roomCategory = NewRoomCategoryUseCase(repos, logger)
	useCases.contract = NewContractUseCase(repos, logger)
	useCases.auth = NewAuthUseCase(repos, logger)

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
