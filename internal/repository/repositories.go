package repository

import (
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/infra/cache"

	"gorm.io/gorm"
)

type repositories struct {
	user               repository.UserRepository
	room               repository.RoomRepository
	roomCategory       repository.RoomCategoryRepository
	contract           repository.ContractRepository
	auth               repository.AuthRepository
	dashboard          repository.DashboardRepository
	email              repository.EmailRepository
	otpCode            repository.OtpCodeRepository
	amenities          repository.AmenitiesRepository
	maintenanceHistory repository.MaintenanceHistoryRepository
	payment            repository.PaymentRepository
}

func NewRepositories(db *gorm.DB, redisClient *cache.RedisClient) repository.Repositories {
	repos := &repositories{}

	repos.user = NewUserRepository(db, redisClient)
	repos.room = NewRoomRepository(db)
	repos.roomCategory = NewRoomCategoryRepository(db)
	repos.contract = NewContractRepository(db)
	repos.auth = NewAuthRepository(db, redisClient)
	repos.dashboard = NewDashboardRepository(db)
	repos.email = NewEmailRepository(db)
	repos.otpCode = NewOtpCodeRepository(db)
	repos.amenities = NewAmenitiesRepository(db)
	repos.maintenanceHistory = NewMaintenanceHistoryRepository(db)
	repos.payment = NewPaymentRepository(db)

	return repos
}

func (r *repositories) User() repository.UserRepository {
	return r.user
}

func (r *repositories) Room() repository.RoomRepository {
	return r.room
}

func (r *repositories) RoomCategory() repository.RoomCategoryRepository {
	return r.roomCategory
}

func (r *repositories) Contract() repository.ContractRepository {
	return r.contract
}

func (r *repositories) Auth() repository.AuthRepository {
	return r.auth
}

func (r *repositories) Dashboard() repository.DashboardRepository {
	return r.dashboard
}

func (r *repositories) OtpCode() repository.OtpCodeRepository {
	return r.otpCode
}

func (r *repositories) Email() repository.EmailRepository {
	return r.email
}

func (r *repositories) Amenities() repository.AmenitiesRepository {
	return r.amenities
}

func (r *repositories) MaintenanceHistory() repository.MaintenanceHistoryRepository {
	return r.maintenanceHistory
}

func (r *repositories) Payment() repository.PaymentRepository {
	return r.payment
}
