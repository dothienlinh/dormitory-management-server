package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type RoomCategoryUseCase interface {
	CreateRoomCategory(ctx context.Context, category *entity.CreateRoomCategory) response.StatusResponse

	GetRoomCategoryByID(ctx context.Context, id uint64) response.StatusResponse

	GetListRoomCategories(ctx context.Context, filter *entity.RoomCategoryFilter) response.StatusResponse

	UpdateRoomCategory(ctx context.Context, id uint64, category *entity.UpdateRoomCategory) response.StatusResponse

	DeleteRoomCategory(ctx context.Context, id uint64) response.StatusResponse
}
