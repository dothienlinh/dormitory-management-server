package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) repository.DashboardRepository {
	return &dashboardRepository{
		db: db,
	}
}

func (r *dashboardRepository) GetStats(ctx context.Context) (*entity.DashboardStats, error) {
	stats := &entity.DashboardStats{
		MonthRevenue: entity.MonthRevenue{
			Month: time.Now().Month().String(),
		},
	}

	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("role = ?", entity.UserRoleStudent).Count(&stats.TotalStudents).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Count(&stats.TotalRooms).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table(entity.Contract{}.TableName()).Count(&stats.TotalContracts).Error; err != nil {
		return nil, err
	}

	// if err := r.db.WithContext(ctx).Table(entity.Contract{}.TableName()).Select("MONTH(created_at) as month, SUM(amount) as amount").Group("MONTH(created_at)").Scan(&stats.MonthRevenue).Error; err != nil {
	// 	return nil, err
	// }

	return stats, nil
}
