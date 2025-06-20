package usecase

import (
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"

	"github.com/hibiken/asynq"
)

type useCases struct {
	user               usecase.UserUseCase
	room               usecase.RoomUseCase
	roomCategory       usecase.RoomCategoryUseCase
	contract           usecase.ContractUseCase
	auth               usecase.AuthUseCase
	dashboard          usecase.DashboardUseCase
	email              usecase.EmailUseCase
	amenities          usecase.AmenitiesUseCase
	maintenanceHistory usecase.MaintenanceHistoryUsecase
}

func NewUseCases(repos repository.Repositories, logger logger.Logger, asynqClient *asynq.Client) usecase.UseCases {
	useCases := &useCases{}

	useCases.user = NewUserUseCase(repos, logger)
	useCases.room = NewRoomUseCase(repos, logger)
	useCases.roomCategory = NewRoomCategoryUseCase(repos, logger)
	useCases.contract = NewContractUseCase(repos, logger)
	useCases.auth = NewAuthUseCase(repos, logger, asynqClient)
	useCases.dashboard = NewDashboardUseCase(repos, logger)
	useCases.email = NewEmailUseCase(repos, logger, asynqClient)
	useCases.amenities = NewAmenitiesUseCase(repos, logger)
	useCases.maintenanceHistory = NewMaintenanceHistoryUsecase(repos, logger)

	return useCases
}

func (uc *useCases) User() usecase.UserUseCase {
	return uc.user
}

func (uc *useCases) Room() usecase.RoomUseCase {
	return uc.room
}

func (uc *useCases) RoomCategory() usecase.RoomCategoryUseCase {
	return uc.roomCategory
}

func (uc *useCases) Contract() usecase.ContractUseCase {
	return uc.contract
}

func (uc *useCases) Auth() usecase.AuthUseCase {
	return uc.auth
}

func (uc *useCases) Dashboard() usecase.DashboardUseCase {
	return uc.dashboard
}

func (uc *useCases) Email() usecase.EmailUseCase {
	return uc.email
}

func (uc *useCases) Amenities() usecase.AmenitiesUseCase {
	return uc.amenities
}

func (uc *useCases) MaintenanceHistory() usecase.MaintenanceHistoryUsecase {
	return uc.maintenanceHistory
}
