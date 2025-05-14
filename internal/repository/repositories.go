package repository

import (
	"dormitory_management/internal/domain/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// repositories implements the repository.Repositories interface
type repositories struct {
	db           *gorm.DB
	redisClient  *redis.Client
	user         repository.UserRepository
	room         repository.RoomRepository
	roomCategory repository.RoomCategoryRepository
	contract     repository.ContractRepository
	auth         repository.AuthRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *gorm.DB, redisClient *redis.Client) repository.Repositories {
	repos := &repositories{
		db:          db,
		redisClient: redisClient,
	}

	repos.user = NewUserRepository(db)
	repos.room = NewRoomRepository(db)
	repos.roomCategory = NewRoomCategoryRepository(db)
	repos.contract = NewContractRepository(db)
	repos.auth = NewAuthRepository(db, redisClient)

	return repos
}

// User returns the user repository
func (r *repositories) User() repository.UserRepository {
	return r.user
}

// Room returns the room repository
func (r *repositories) Room() repository.RoomRepository {
	return r.room
}

// RoomCategory returns the room category repository
func (r *repositories) RoomCategory() repository.RoomCategoryRepository {
	return r.roomCategory
}

// Contract returns the contract repository
func (r *repositories) Contract() repository.ContractRepository {
	return r.contract
}

// Auth returns the auth repository
func (r *repositories) Auth() repository.AuthRepository {
	return r.auth
}
