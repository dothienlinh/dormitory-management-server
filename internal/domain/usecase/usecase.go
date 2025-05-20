package usecase

// UseCases is a collection of all use cases
type UseCases interface {
	User() UserUseCase
	Room() RoomUseCase
	RoomCategory() RoomCategoryUseCase
	Contract() ContractUseCase
	Auth() AuthUseCase
	Dashboard() DashboardUseCase
}
