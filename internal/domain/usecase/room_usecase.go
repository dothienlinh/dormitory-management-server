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

	// Student room methods
	GetStudentRoomDetails(ctx context.Context, userID uint64) response.StatusResponse
	GetRoomStats(ctx context.Context, userID uint64) response.StatusResponse

	// Room issues
	GetRoomIssues(ctx context.Context, userID uint64, page, limit int, status, category string) response.StatusResponse
	CreateRoomIssue(ctx context.Context, userID uint64, title, description, category, priority string) response.StatusResponse
	GetRoomIssueDetails(ctx context.Context, issueID uint64, userID uint64) response.StatusResponse

	// Room bills
	GetRoomBills(ctx context.Context, userID uint64, page, limit, year, month int, status string) response.StatusResponse

	// Room rules
	GetRoomRules(ctx context.Context) response.StatusResponse

	// Cleaning schedule
	GetCleaningSchedule(ctx context.Context, userID uint64) response.StatusResponse

	// Roommates
	GetRoommates(ctx context.Context, userID uint64) response.StatusResponse
}
