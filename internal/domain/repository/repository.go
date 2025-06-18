package repository

type Repositories interface {
	User() UserRepository
	Room() RoomRepository
	RoomCategory() RoomCategoryRepository
	Contract() ContractRepository
	Auth() AuthRepository
	Dashboard() DashboardRepository
	Email() EmailRepository
	OtpCode() OtpCodeRepository
	Facilities() FacilitiesRepository
}
