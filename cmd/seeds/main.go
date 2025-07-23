package main

import (
	"context"
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/infra/database"
	"dormitory_management/pkg/logger"
	"flag"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func main() {
	context := context.Background()

	cfg := config.LoadConfig()

	log := logger.NewLogger(cfg.LogLevel)

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	tableFlags := flag.String("tables", "", "Comma-separated list of tables to seed (e.g., 'amenities,room_categories,notifications,events,service_requests,room_rules,cleaning_schedules')")
	flag.Parse()

	tableStrings := *tableFlags
	if tableStrings == "" {
		log.Fatal("No tables specified for seeding. Use the -tables flag to specify which tables to seed.", fmt.Errorf("usage: %s -tables=amenities,room_categories,notifications,events,service_requests,room_rules,cleaning_schedules", flag.CommandLine.Name()))
	}

	log.Info("Seeding tables...")

	switch tableStrings {
	case entity.Amenity{}.TableName():
		if err := seedAmenities(context, db); err != nil {
			log.Fatal("Failed to seed amenities", err)
		}
	case entity.RoomCategory{}.TableName():
		if err := seedRoomCategories(context, db); err != nil {
			log.Fatal("Failed to seed room categories", err)
		}
	case entity.Notification{}.TableName():
		if err := seedNotifications(context, db); err != nil {
			log.Fatal("Failed to seed notifications", err)
		}
	case entity.Event{}.TableName():
		if err := seedEvents(context, db); err != nil {
			log.Fatal("Failed to seed events", err)
		}
	case entity.ServiceRequest{}.TableName():
		if err := seedServiceRequests(context, db); err != nil {
			log.Fatal("Failed to seed service requests", err)
		}
	case entity.RoomRule{}.TableName():
		if err := seedRoomRules(context, db); err != nil {
			log.Fatal("Failed to seed room rules", err)
		}
	case entity.CleaningSchedule{}.TableName():
		if err := seedCleaningSchedules(context, db); err != nil {
			log.Fatal("Failed to seed cleaning schedules", err)
		}
	case "all":
		if err := seedAmenities(context, db); err != nil {
			log.Fatal("Failed to seed amenities", err)
		}
		if err := seedRoomCategories(context, db); err != nil {
			log.Fatal("Failed to seed room categories", err)
		}
		if err := seedNotifications(context, db); err != nil {
			log.Fatal("Failed to seed notifications", err)
		}
		if err := seedEvents(context, db); err != nil {
			log.Fatal("Failed to seed events", err)
		}
		if err := seedServiceRequests(context, db); err != nil {
			log.Fatal("Failed to seed service requests", err)
		}
		if err := seedRoomRules(context, db); err != nil {
			log.Fatal("Failed to seed room rules", err)
		}
		if err := seedCleaningSchedules(context, db); err != nil {
			log.Fatal("Failed to seed cleaning schedules", err)
		}
		log.Info("All tables seeded successfully")
		return
	default:
		log.Fatal("Invalid table specified. Use -tables=amenities,room_categories,notifications,events,service_requests,room_rules,cleaning_schedules to specify which tables to seed.", fmt.Errorf("usage: %s -tables=amenities,room_categories,notifications,events,service_requests,room_rules,cleaning_schedules", flag.CommandLine.Name()))
	}

	log.Info("Seeding completed successfully")
}

func seedAmenities(ctx context.Context, db *gorm.DB) error {
	amenities := []entity.Amenity{
		{Name: "Wifi"},
		{Name: "Điều hoà"},
		{Name: "Nóng lạnh"},
		{Name: "Tủ đồ"},
		{Name: "Quạt trần"},
		{Name: "Bàn ghế"},
	}

	return db.WithContext(ctx).Table(entity.Amenity{}.TableName()).Create(&amenities).Error
}

func seedRoomCategories(ctx context.Context, db *gorm.DB) error {
	roomCategories := []entity.RoomCategory{
		{Name: "Phòng 4", Description: "Phòng dành cho 4 người", Capacity: 4, Price: 1800000, Acreage: 20},
		{Name: "Phòng 6", Description: "Phòng dành cho 6 người", Capacity: 6, Price: 1600000, Acreage: 25},
		{Name: "Phòng 8", Description: "Phòng dành cho 8 người", Capacity: 8, Price: 1400000, Acreage: 30},
	}

	return db.WithContext(ctx).Table(entity.RoomCategory{}.TableName()).Create(&roomCategories).Error
}

func seedNotifications(ctx context.Context, db *gorm.DB) error {
	// Get a sample user ID (assuming there's at least one user)
	var userID uint64
	if err := db.WithContext(ctx).Model(&entity.User{}).Select("id").Where("role = ?", entity.UserRoleStudent).First(&userID).Error; err != nil {
		return fmt.Errorf("no student user found to create notifications: %w", err)
	}

	notifications := []entity.Notification{
		{
			Title:    "Thông báo thanh toán tiền phòng tháng 12",
			Content:  "Vui lòng thanh toán tiền phòng tháng 12 trước ngày 25/12/2023. Số tiền: 2,500,000 VND",
			Type:     entity.NotificationTypePayment,
			Priority: entity.NotificationPriorityHigh,
			UserID:   userID,
			IsRead:   false,
		},
		{
			Title:    "Bảo trì hệ thống điện",
			Content:  "Hệ thống điện tòa A sẽ được bảo trì vào ngày 20/12/2023 từ 8:00-12:00",
			Type:     entity.NotificationTypeMaintenance,
			Priority: entity.NotificationPriorityMedium,
			UserID:   userID,
			IsRead:   true,
		},
		{
			Title:    "Lễ hội cuối năm KTX",
			Content:  "Tham gia lễ hội cuối năm KTX vào ngày 23/12/2023 tại sân chính",
			Type:     entity.NotificationTypeEvent,
			Priority: entity.NotificationPriorityLow,
			UserID:   userID,
			IsRead:   false,
		},
		{
			Title:    "Yêu cầu sửa chữa đã được xử lý",
			Content:  "Yêu cầu sửa chữa đèn phòng A304 đã được hoàn thành",
			Type:     entity.NotificationTypeService,
			Priority: entity.NotificationPriorityMedium,
			UserID:   userID,
			IsRead:   false,
		},
		{
			Title:    "Thông báo quy định mới",
			Content:  "KTX có quy định mới về giờ giấc sinh hoạt, vui lòng xem chi tiết",
			Type:     entity.NotificationTypeAnnouncement,
			Priority: entity.NotificationPriorityMedium,
			UserID:   userID,
			IsRead:   true,
		},
	}

	return db.WithContext(ctx).Create(&notifications).Error
}

func seedEvents(ctx context.Context, db *gorm.DB) error {
	now := time.Now()
	events := []entity.Event{
		{
			Title:       "Hạn nộp tiền phòng tháng 12",
			Description: "Hạn cuối cùng nộp tiền phòng tháng 12",
			EventDate:   now.AddDate(0, 0, 10), // 10 ngày từ bây giờ
			StartTime:   "08:00",
			EndTime:     "17:00",
			Location:    "Phòng kế toán - Tầng 1",
			Type:        entity.EventTypeDeadline,
			IsMandatory: true,
			Organizer:   "Ban quản lý KTX",
		},
		{
			Title:       "Bảo trì hệ thống nước",
			Description: "Bảo trì và vệ sinh bồn nước các tòa nhà",
			EventDate:   now.AddDate(0, 0, 5), // 5 ngày từ bây giờ
			StartTime:   "08:00",
			EndTime:     "16:00",
			Location:    "Tất cả các tòa nhà",
			Type:        entity.EventTypeMaintenance,
			IsMandatory: false,
			Organizer:   "Phòng kỹ thuật",
		},
		{
			Title:       "Họp sinh viên đầu tháng",
			Description: "Họp tổng kết tháng và thông báo các hoạt động sắp tới",
			EventDate:   now.AddDate(0, 0, 2), // 2 ngày từ bây giờ
			StartTime:   "19:00",
			EndTime:     "21:00",
			Location:    "Hội trường tầng 1",
			Type:        entity.EventTypeMeeting,
			IsMandatory: true,
			Organizer:   "Ban sinh viên KTX",
		},
		{
			Title:       "Đêm nhạc Giáng sinh",
			Description: "Chương trình văn nghệ chào mừng Giáng sinh",
			EventDate:   now.AddDate(0, 0, 8), // 8 ngày từ bây giờ
			StartTime:   "20:00",
			EndTime:     "22:00",
			Location:    "Sân chính KTX",
			Type:        entity.EventTypeEvent,
			IsMandatory: false,
			Organizer:   "Đoàn thanh niên KTX",
		},
		{
			Title:       "Kiểm tra an toàn PCCC",
			Description: "Kiểm tra hệ thống phòng cháy chữa cháy định kỳ",
			EventDate:   now.AddDate(0, 0, 15), // 15 ngày từ bây giờ
			StartTime:   "08:00",
			EndTime:     "12:00",
			Location:    "Tất cả các tòa nhà",
			Type:        entity.EventTypeMaintenance,
			IsMandatory: false,
			Organizer:   "Phòng an ninh",
		},
	}

	return db.WithContext(ctx).Create(&events).Error
}

func seedServiceRequests(ctx context.Context, db *gorm.DB) error {
	// Get a sample user ID (assuming there's at least one user)
	var userID uint64
	if err := db.WithContext(ctx).Model(&entity.User{}).Select("id").Where("role = ?", entity.UserRoleStudent).First(&userID).Error; err != nil {
		return fmt.Errorf("no student user found to create service requests: %w", err)
	}

	now := time.Now()
	serviceRequests := []entity.ServiceRequest{
		{
			Title:          "Sửa chữa đèn phòng A304",
			Description:    "Đèn trong phòng bị hỏng, không sáng được",
			Category:       "Điện",
			Priority:       entity.ServiceRequestPriorityHigh,
			Status:         entity.ServiceRequestStatusCompleted,
			UserID:         userID,
			CompletionDate: &now,
		},
		{
			Title:       "Thay ổ khóa cửa phòng",
			Description: "Ổ khóa cửa phòng bị kẹt, khó mở",
			Category:    "Cơ khí",
			Priority:    entity.ServiceRequestPriorityMedium,
			Status:      entity.ServiceRequestStatusPending,
			UserID:      userID,
		},
		{
			Title:       "Sửa chữa vòi nước",
			Description: "Vòi nước trong toilet bị rò rỉ",
			Category:    "Nước",
			Priority:    entity.ServiceRequestPriorityMedium,
			Status:      entity.ServiceRequestStatusInProgress,
			UserID:      userID,
		},
		{
			Title:       "Thay bóng đèn hành lang",
			Description: "Bóng đèn hành lang tầng 3 bị cháy",
			Category:    "Điện",
			Priority:    entity.ServiceRequestPriorityLow,
			Status:      entity.ServiceRequestStatusApproved,
			UserID:      userID,
		},
		{
			Title:       "Sửa chữa điều hòa",
			Description: "Điều hòa không hoạt động, không làm lạnh",
			Category:    "Điện",
			Priority:    entity.ServiceRequestPriorityHigh,
			Status:      entity.ServiceRequestStatusPending,
			UserID:      userID,
		},
	}

	return db.WithContext(ctx).Create(&serviceRequests).Error
}

func seedRoomRules(ctx context.Context, db *gorm.DB) error {
	roomRules := []entity.RoomRule{
		{
			Title:       "Giờ giấc tự do",
			Description: "Giờ giấc tự do, không quy định đóng cửa",
			Category:    "general",
			Priority:    1,
			IsActive:    true,
		},
		{
			Title:       "Không mang thú cưng",
			Description: "Không được mang thú cưng vào phòng",
			Category:    "general",
			Priority:    2,
			IsActive:    true,
		},
		{
			Title:       "Giữ gìn vệ sinh",
			Description: "Giữ gìn vệ sinh chung, không xả rác bừa bãi",
			Category:    "hygiene",
			Priority:    1,
			IsActive:    true,
		},
		{
			Title:       "Báo cáo sự cố an ninh",
			Description: "Báo ngay cho bảo vệ khi phát hiện sự cố về an ninh trật tự",
			Category:    "safety",
			Priority:    1,
			IsActive:    true,
		},
		{
			Title:       "Tham gia họp tổ",
			Description: "Tham gia đầy đủ các buổi họp tổ dân phố định kỳ",
			Category:    "behavior",
			Priority:    3,
			IsActive:    true,
		},
	}

	return db.WithContext(ctx).Create(&roomRules).Error
}

func seedCleaningSchedules(ctx context.Context, db *gorm.DB) error {
	cleaningSchedules := []entity.CleaningSchedule{
		{
			RoomID:    nil, // applies to all rooms
			DayOfWeek: "monday",
			StartTime: "08:00",
			EndTime:   "10:00",
			Type:      "regular",
			IsActive:  true,
		},
		{
			RoomID:    nil,
			DayOfWeek: "tuesday",
			StartTime: "10:00",
			EndTime:   "12:00",
			Type:      "regular",
			IsActive:  true,
		},
		{
			RoomID:    nil,
			DayOfWeek: "wednesday",
			StartTime: "08:00",
			EndTime:   "10:00",
			Type:      "regular",
			IsActive:  true,
		},
		{
			RoomID:    nil,
			DayOfWeek: "thursday",
			StartTime: "10:00",
			EndTime:   "12:00",
			Type:      "regular",
			IsActive:  true,
		},
		{
			RoomID:    nil,
			DayOfWeek: "friday",
			StartTime: "08:00",
			EndTime:   "10:00",
			Type:      "regular",
			IsActive:  true,
		},
		{
			RoomID:    nil,
			DayOfWeek: "saturday",
			StartTime: "09:00",
			EndTime:   "12:00",
			Type:      "deep",
			IsActive:  true,
		},
	}

	return db.WithContext(ctx).Create(&cleaningSchedules).Error
}
