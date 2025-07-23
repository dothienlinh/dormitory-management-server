package usecase

type UseCases interface {
	User() UserUseCase
	Room() RoomUseCase
	RoomCategory() RoomCategoryUseCase
	Contract() ContractUseCase
	Auth() AuthUseCase
	Dashboard() DashboardUseCase
	Email() EmailUseCase
	Amenities() AmenitiesUseCase
	MaintenanceHistory() MaintenanceHistoryUsecase
	Payment() PaymentUseCase
	PaymentHistory() PaymentHistoryUseCase
	Bill() BillUseCase
}
