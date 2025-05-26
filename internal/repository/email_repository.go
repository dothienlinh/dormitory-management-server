package repository

import (
	"gorm.io/gorm"
)

type emailRepository struct {
	db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) *emailRepository {
	return &emailRepository{db: db}
}
