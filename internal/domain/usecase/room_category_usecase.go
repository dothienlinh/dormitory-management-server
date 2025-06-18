package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

// RoomCategoryUseCase defines the interface for room category business logic
type RoomCategoryUseCase interface {
	// CreateRoomCategory creates a new room category
	CreateRoomCategory(ctx context.Context, category *entity.CreateRoomCategory) response.StatusResponse

	// GetRoomCategoryByID retrieves a room category by ID
	GetRoomCategoryByID(ctx context.Context, id uint64) response.StatusResponse

	// GetListRoomCategories retrieves room categories based on filter
	GetListRoomCategories(ctx context.Context, filter *entity.RoomCategoryFilter) response.StatusResponse

	// UpdateRoomCategory updates a room category
	UpdateRoomCategory(ctx context.Context, id uint64, category *entity.UpdateRoomCategory) response.StatusResponse

	// DeleteRoomCategory deletes a room category
	DeleteRoomCategory(ctx context.Context, id uint64) response.StatusResponse
}
