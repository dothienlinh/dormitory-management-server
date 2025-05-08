package db

import (
	"dormitory_management/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBClient() *gorm.DB {
	dsn := os.Getenv("DB_CONN_STR")
	dbClient, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln("Failed to connect to database", err)
	}

	if err := dbClient.Exec(CreateEnumTypes).Error; err != nil {
		log.Fatalln("Failed to create enum types", err)
	}

	if err := dbClient.AutoMigrate(&models.User{}, &models.Room{}, &models.Amenities{}, &models.RoomCategory{}); err != nil {
		log.Fatalln("Failed to migrate database", err)
	}

	return dbClient
}
