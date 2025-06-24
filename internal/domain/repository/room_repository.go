package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type RoomRepository interface {
	Create(ctx context.Context, room *entity.CreateRoom) error

	GetByID(ctx context.Context, id uint64) (*entity.Room, error)

	List(ctx context.Context, filter *entity.RoomFilter) ([]entity.ListRooms, int64, error)

	Update(ctx context.Context, room *entity.Room, amenityIDs []uint64) error

	Delete(ctx context.Context, id uint64) error

	AmountStudentsInRoom(ctx context.Context, roomID uint64, amount *int64) error
}
