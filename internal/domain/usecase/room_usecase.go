package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

// RoomUseCase defines the interface for room business logic
type RoomUseCase interface {
	// CreateRoom creates a new room
	CreateRoom(ctx context.Context, room *entity.Room) response.StatusResponse

	// GetRoomByID retrieves a room by ID
	GetRoomByID(ctx context.Context, id uint) response.StatusResponse

	// GetListRooms retrieves rooms based on filter
	GetListRooms(ctx context.Context, filter *entity.RoomFilter) response.StatusResponse

	// UpdateRoom updates a room
	UpdateRoom(ctx context.Context, id uint, room *entity.Room) response.StatusResponse

	// DeleteRoom deletes a room
	DeleteRoom(ctx context.Context, id uint) response.StatusResponse
}

// RoomCategoryUseCase defines the interface for room category business logic
type RoomCategoryUseCase interface {
	// CreateRoomCategory creates a new room category
	CreateRoomCategory(ctx context.Context, category *entity.RoomCategory) response.StatusResponse

	// GetRoomCategoryByID retrieves a room category by ID
	GetRoomCategoryByID(ctx context.Context, id uint) response.StatusResponse

	// GetListRoomCategories retrieves room categories based on filter
	GetListRoomCategories(ctx context.Context, filter *entity.RoomCategoryFilter) response.StatusResponse

	// UpdateRoomCategory updates a room category
	UpdateRoomCategory(ctx context.Context, id uint, category *entity.RoomCategory) response.StatusResponse

	// DeleteRoomCategory deletes a room category
	DeleteRoomCategory(ctx context.Context, id uint) response.StatusResponse
}
