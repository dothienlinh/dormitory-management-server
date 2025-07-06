package main

import (
	"context"
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/infra/database"
	"dormitory_management/pkg/logger"
	"flag"
	"fmt"

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

	tableFlags := flag.String("tables", "", "Comma-separated list of tables to seed (e.g., 'amenities,room_categories')")
	flag.Parse()

	tableStrings := *tableFlags
	if tableStrings == "" {
		log.Fatal("No tables specified for seeding. Use the -tables flag to specify which tables to seed.", fmt.Errorf("usage: %s -tables=amenities,room_categories", flag.CommandLine.Name()))
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
	default:
		log.Fatal("Invalid table specified. Use -tables=amenities,room_categories to specify which tables to seed.", fmt.Errorf("usage: %s -tables=amenities,room_categories", flag.CommandLine.Name()))
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
