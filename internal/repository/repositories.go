package repository

import (
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/infra/cache"

	"gorm.io/gorm"
)

// repositories implements the repository.Repositories interface
type repositories struct {
	user         repository.UserRepository
	room         repository.RoomRepository
	roomCategory repository.RoomCategoryRepository
	contract     repository.ContractRepository
	auth         repository.AuthRepository
	dashboard    repository.DashboardRepository
	email        repository.EmailRepository
	otpCode      repository.OtpCodeRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *gorm.DB, redisClient *cache.RedisClient) repository.Repositories {
	repos := &repositories{}

	repos.user = NewUserRepository(db)
	repos.room = NewRoomRepository(db)
	repos.roomCategory = NewRoomCategoryRepository(db)
	repos.contract = NewContractRepository(db)
	repos.auth = NewAuthRepository(db, redisClient)
	repos.dashboard = NewDashboardRepository(db)
	repos.email = NewEmailRepository(db)
	repos.otpCode = NewOtpCodeRepository(db)

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

// Dashboard returns the dashboard repository
func (r *repositories) Dashboard() repository.DashboardRepository {
	return r.dashboard
}

// OtpCode returns the OTP code repository
func (r *repositories) OtpCode() repository.OtpCodeRepository {
	return r.otpCode
}

// Email returns the email repository
func (r *repositories) Email() repository.EmailRepository {
	return r.email
}
