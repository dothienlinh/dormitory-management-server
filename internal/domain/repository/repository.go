package repository

// Repositories is a collection of all repositories
type Repositories interface {
	User() UserRepository
	Room() RoomRepository
	RoomCategory() RoomCategoryRepository
	Contract() ContractRepository
	Auth() AuthRepository
	Dashboard() DashboardRepository
	Email() EmailRepository
	OtpCode() OtpCodeRepository
}
