package handler

import (
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	User         *UserHandler
	Room         *RoomHandler
	RoomCategory *RoomCategoryHandler
	Contract     *ContractHandler
	Auth         *AuthHandler
	Dashboard    *DashboardHandler
	Email        *EmailHandler
}

// NewHandlers creates a new Handlers instance
func NewHandlers(useCases usecase.UseCases, logger logger.Logger) *Handlers {
	return &Handlers{
		User:         NewUserHandler(useCases, logger),
		Room:         NewRoomHandler(useCases, logger),
		RoomCategory: NewRoomCategoryHandler(useCases, logger),
		Contract:     NewContractHandler(useCases, logger),
		Auth:         NewAuthHandler(useCases, logger),
		Dashboard:    NewDashboardHandler(useCases, logger),
		Email:        NewEmailHandler(useCases, logger),
	}
}
