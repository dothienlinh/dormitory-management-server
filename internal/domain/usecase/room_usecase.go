package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

// RoomUseCase defines the interface for room business logic
type RoomUseCase interface {
	// CreateRoom creates a new room
	CreateRoom(ctx context.Context, room *entity.CreateRoom) response.StatusResponse

	// GetRoomByID retrieves a room by ID
	GetRoomByID(ctx context.Context, id uint64) response.StatusResponse

	// GetListRooms retrieves rooms based on filter
	GetListRooms(ctx context.Context, filter *entity.RoomFilter) response.StatusResponse

	// UpdateRoom updates a room
	UpdateRoom(ctx context.Context, id uint64, room *entity.UpdateRoom) response.StatusResponse

	// DeleteRoom deletes a room
	DeleteRoom(ctx context.Context, id uint64) response.StatusResponse
}
