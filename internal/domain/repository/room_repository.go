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

	// Student room methods
	GetStudentRoomDetails(ctx context.Context, userID uint64) (*entity.StudentRoomDetails, error)
	GetRoomStats(ctx context.Context, userID uint64) (*entity.RoomStats, error)

	// Room issues
	GetRoomIssues(ctx context.Context, userID uint64, params *RoomIssueParams) ([]*entity.RoomIssue, int, error)
	CreateRoomIssue(ctx context.Context, issue *entity.RoomIssue) error
	GetRoomIssueByID(ctx context.Context, issueID uint64, userID uint64) (*entity.RoomIssue, error)

	// Room bills
	GetRoomBills(ctx context.Context, userID uint64, params *RoomBillParams) ([]*entity.Payment, int, error)

	// Room rules
	GetRoomRules(ctx context.Context) ([]*entity.RoomRule, error)

	// Cleaning schedule
	GetCleaningSchedule(ctx context.Context, roomID uint64) ([]*entity.CleaningSchedule, error)

	// Roommates
	GetRoommates(ctx context.Context, userID uint64) ([]*entity.User, error)
}

type RoomIssueParams struct {
	Page     int
	Limit    int
	Status   string
	Category string
}

type RoomBillParams struct {
	Page   int
	Limit  int
	Year   int
	Month  int
	Status string
}
