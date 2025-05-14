package entity

import (
	"time"

	"gorm.io/gorm"
)

// Base model for all entities
type Base struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Pagination represents pagination parameters for queries
type Pagination struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
	Total int `json:"total"`
}

// Parse sets default values for pagination if not provided
func (p *Pagination) Parse() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 10
	}
}

// GetOffset returns the offset for pagination queries
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
