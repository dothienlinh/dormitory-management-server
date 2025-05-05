package db

import (
	"dormitory_management/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_CONN_STR")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Failed to connect to database", err)
	}

	if err := db.Exec(CreateEnumTypes).Error; err != nil {
		log.Fatalln("Failed to create enum types", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Room{}); err != nil {
		log.Fatalln("Failed to migrate database", err)
	}

	DB = db
}
