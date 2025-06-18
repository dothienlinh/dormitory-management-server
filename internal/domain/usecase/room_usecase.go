package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type RoomUseCase interface {
	CreateRoom(ctx context.Context, room *entity.CreateRoom) response.StatusResponse

	GetRoomByID(ctx context.Context, id uint64) response.StatusResponse

	GetListRooms(ctx context.Context, filter *entity.RoomFilter) response.StatusResponse

	UpdateRoom(ctx context.Context, id uint64, room *entity.UpdateRoom) response.StatusResponse

	DeleteRoom(ctx context.Context, id uint64) response.StatusResponse
}
