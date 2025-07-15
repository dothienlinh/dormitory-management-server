package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/helper"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &roomRepository{
		db: db,
	}
}

func (r *roomRepository) Create(ctx context.Context, payload *entity.CreateRoom) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		roomEntity := &entity.Room{
			RoomNumber:     payload.RoomNumber,
			Status:         payload.Status,
			RoomCategoryID: payload.RoomCategoryID,
		}

		if err := tx.WithContext(ctx).Create(roomEntity).Error; err != nil {
			return fmt.Errorf("Create to create room: %w", err)
		}

		if len(payload.AmenityIDs) > 0 {
			uniqueAmenityIDs := helper.RemoveDuplicateUint64(payload.AmenityIDs)
			var existingCount int64
			if err := tx.WithContext(ctx).
				Table(entity.RoomAmenities{}.TableName()).
				Where("room_id = ? AND amenity_id IN ?", roomEntity.ID, uniqueAmenityIDs).
				Count(&existingCount).Error; err != nil {
				return fmt.Errorf("failed to check existing room amenities: %w", err)
			}

			if existingCount > 0 {
				return fmt.Errorf("some amenities are already associated with this room")
			}

			if err := r.createRoomAmenities(ctx, tx, roomEntity.ID, uniqueAmenityIDs); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *roomRepository) createRoomAmenities(ctx context.Context, tx *gorm.DB, roomID uint64, amenityIDs []uint64) error {
	var amenityCount int64
	if err := tx.WithContext(ctx).
		Table(entity.Amenity{}.TableName()).
		Where("id IN ?", amenityIDs).
		Count(&amenityCount).Error; err != nil {
		return fmt.Errorf("failed to validate amenities: %w", err)
	}

	if int64(len(amenityIDs)) != amenityCount {
		return fmt.Errorf("some amenity IDs do not exist")
	}

	roomAmenities := make([]entity.RoomAmenities, 0, len(amenityIDs))
	for _, amenityID := range amenityIDs {
		roomAmenities = append(roomAmenities, entity.RoomAmenities{
			RoomID:    roomID,
			AmenityID: amenityID,
		})
	}

	if err := tx.WithContext(ctx).Create(&roomAmenities).Error; err != nil {
		return fmt.Errorf("failed to create room amenities: %w", err)
	}

	return nil
}

func (r *roomRepository) GetByID(ctx context.Context, id uint64) (*entity.Room, error) {
	room := entity.Room{}
	if err := r.db.WithContext(ctx).Table("rooms r").
		Preload("Users").
		Preload("RoomCategory").
		Preload("RoomAmenities.Amenity").
		Preload("MaintenanceHistory").
		Select("r.*").
		Where("r.id = ?", id).
		First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}
	return &room, nil
}

func (r *roomRepository) List(ctx context.Context, filter *entity.RoomFilter) ([]entity.ListRooms, int64, error) {
	filter.Parse()
	whereClause, values := filter.Build()

	var results []entity.ListRooms
	var total int64

	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Where(whereClause, values...).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	if err := r.db.WithContext(ctx).Table("rooms r").Where(whereClause, values...).
		Select("r.*, COUNT(u.id) as user_count").
		Joins("LEFT JOIN users u ON r.id = u.room_id").
		Preload("RoomCategory").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Group("r.id").
		Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get rooms: %w", err)
	}

	return results, total, nil
}

func (r *roomRepository) Update(ctx context.Context, room *entity.Room, amenityIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Table(room.TableName()).Save(room).Error; err != nil {
			return fmt.Errorf("failed to update room: %w", err)
		}

		roomAmenities := entity.RoomAmenities{
			RoomID: room.ID,
		}

		if err := tx.WithContext(ctx).Table(roomAmenities.TableName()).Unscoped().Where("room_id = ?", room.ID).Delete(&roomAmenities).Error; err != nil {
			return fmt.Errorf("failed to delete existing room amenities: %w", err)
		}

		if len(amenityIDs) > 0 {
			uniqueAmenityIDs := helper.RemoveDuplicateUint64(amenityIDs)
			if err := r.createRoomAmenities(ctx, tx, room.ID, uniqueAmenityIDs); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *roomRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Delete(&entity.Room{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}

func (r *roomRepository) AmountStudentsInRoom(ctx context.Context, roomID uint64, amount *int64) error {
	return r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("room_id = ?", roomID).Count(amount).Error
}

// GetStudentRoomDetails retrieves detailed room information for a student
func (r *roomRepository) GetStudentRoomDetails(ctx context.Context, userID uint64) (*entity.StudentRoomDetails, error) {
	// Get user's room information
	var user entity.User
	if err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("Room.RoomCategory").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.RoomID == nil {
		return nil, errors.New("user has no room assigned")
	}

	// Get room details with all relations
	var room entity.Room
	if err := r.db.WithContext(ctx).
		Preload("RoomCategory").
		Preload("RoomAmenities.Amenity").
		Where("id = ?", *user.RoomID).
		First(&room).Error; err != nil {
		return nil, fmt.Errorf("failed to get room details: %w", err)
	}

	// Get roommates
	var roommates []*entity.User
	if err := r.db.WithContext(ctx).
		Where("room_id = ? AND id != ?", *user.RoomID, userID).
		Find(&roommates).Error; err != nil {
		return nil, fmt.Errorf("failed to get roommates: %w", err)
	}

	// Get recent issues
	var recentIssues []*entity.RoomIssue
	if err := r.db.WithContext(ctx).
		Preload("Reporter").
		Where("room_id = ?", *user.RoomID).
		Order("created_at DESC").
		Limit(5).
		Find(&recentIssues).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent issues: %w", err)
	}

	// Get cleaning schedule
	var cleaningSchedule []*entity.CleaningSchedule
	if err := r.db.WithContext(ctx).
		Where("room_id = ? OR room_id IS NULL", *user.RoomID).
		Where("is_active = ?", true).
		Find(&cleaningSchedule).Error; err != nil {
		return nil, fmt.Errorf("failed to get cleaning schedule: %w", err)
	}

	// Get bill history
	var billHistory []*entity.Payment
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(5).
		Find(&billHistory).Error; err != nil {
		return nil, fmt.Errorf("failed to get bill history: %w", err)
	}

	// Convert amenities to facilities
	var facilities []*entity.Amenity
	for _, roomAmenity := range room.RoomAmenities {
		if roomAmenity.Amenity != nil {
			facilities = append(facilities, roomAmenity.Amenity)
		}
	}

	// Build response
	details := &entity.StudentRoomDetails{
		ID:               fmt.Sprintf("%d", room.ID),
		Name:             room.RoomNumber,
		BuildingName:     "Tòa nhà A", // You might want to add building info to Room entity
		BuildingAddress:  "123 Đường ABC, Quận 1, TP.HCM",
		Floor:            3, // You might want to extract this from room number
		RoomType:         room.RoomCategory.Name,
		Capacity:         room.RoomCategory.Capacity,
		Area:             float64(room.RoomCategory.Acreage),
		MonthlyPrice:     float64(room.RoomCategory.Price),
		CurrentOccupants: len(roommates) + 1, // +1 for current user
		Status:           string(room.Status),
		Facilities:       facilities,
		Residents:        append(roommates, &user),
		RecentIssues:     recentIssues,
		CleaningSchedule: cleaningSchedule,
		BillHistory:      billHistory,
		Room:             &room,
	}

	return details, nil
}

// GetRoomStats retrieves room statistics for a student
func (r *roomRepository) GetRoomStats(ctx context.Context, userID uint64) (*entity.RoomStats, error) {
	// Get user's room
	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.RoomID == nil {
		return nil, errors.New("user has no room assigned")
	}

	stats := &entity.RoomStats{}

	// Count facilities
	var totalFacilities int64
	r.db.WithContext(ctx).Model(&entity.RoomAmenities{}).Where("room_id = ?", *user.RoomID).Count(&totalFacilities)
	stats.TotalFacilities = int(totalFacilities)
	stats.WorkingFacilities = int(totalFacilities) // Assuming all are working for now

	// Count issues
	var pendingIssues, resolvedIssues int64
	r.db.WithContext(ctx).Model(&entity.RoomIssue{}).Where("room_id = ? AND status = ?", *user.RoomID, "pending").Count(&pendingIssues)
	r.db.WithContext(ctx).Model(&entity.RoomIssue{}).Where("room_id = ? AND status = ?", *user.RoomID, "resolved").Count(&resolvedIssues)
	stats.PendingIssues = int(pendingIssues)
	stats.ResolvedIssues = int(resolvedIssues)

	// Count current occupancy
	var currentOccupancy int64
	r.db.WithContext(ctx).Model(&entity.User{}).Where("room_id = ?", *user.RoomID).Count(&currentOccupancy)
	stats.CurrentOccupancy = int(currentOccupancy)

	// Get next cleaning (simplified)
	stats.NextCleaning = "Thứ 2, 08:00"

	return stats, nil
}

// GetRoomIssues retrieves room issues with pagination
func (r *roomRepository) GetRoomIssues(ctx context.Context, userID uint64, params *repository.RoomIssueParams) ([]*entity.RoomIssue, int, error) {
	// Get user's room
	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get user: %w", err)
	}

	if user.RoomID == nil {
		return nil, 0, errors.New("user has no room assigned")
	}

	if params == nil {
		params = &repository.RoomIssueParams{Page: 1, Limit: 10}
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var issues []*entity.RoomIssue
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.RoomIssue{}).Where("room_id = ?", *user.RoomID)

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	if err := query.Preload("Reporter").Preload("AssignedUser").
		Offset(offset).Limit(params.Limit).
		Order("created_at DESC").
		Find(&issues).Error; err != nil {
		return nil, 0, err
	}

	return issues, int(total), nil
}

// CreateRoomIssue creates a new room issue
func (r *roomRepository) CreateRoomIssue(ctx context.Context, issue *entity.RoomIssue) error {
	return r.db.WithContext(ctx).Create(issue).Error
}

// GetRoomIssueByID retrieves a room issue by ID
func (r *roomRepository) GetRoomIssueByID(ctx context.Context, issueID uint64, userID uint64) (*entity.RoomIssue, error) {
	// Get user's room first to verify access
	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.RoomID == nil {
		return nil, errors.New("user has no room assigned")
	}

	var issue entity.RoomIssue
	if err := r.db.WithContext(ctx).
		Preload("Reporter").
		Preload("AssignedUser").
		Where("id = ? AND room_id = ?", issueID, *user.RoomID).
		First(&issue).Error; err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return &issue, nil
}

// GetRoomBills retrieves room bills for a user
func (r *roomRepository) GetRoomBills(ctx context.Context, userID uint64, params *repository.RoomBillParams) ([]*entity.Payment, int, error) {
	if params == nil {
		params = &repository.RoomBillParams{Page: 1, Limit: 10}
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var bills []*entity.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Payment{}).Where("user_id = ?", userID)

	// Apply filters
	if params.Year > 0 {
		query = query.Where("EXTRACT(YEAR FROM created_at) = ?", params.Year)
	}
	if params.Month > 0 {
		query = query.Where("EXTRACT(MONTH FROM created_at) = ?", params.Month)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).
		Order("created_at DESC").
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, int(total), nil
}

// GetRoomRules retrieves all active room rules
func (r *roomRepository) GetRoomRules(ctx context.Context) ([]*entity.RoomRule, error) {
	var rules []*entity.RoomRule
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("priority ASC, created_at ASC").
		Find(&rules).Error; err != nil {
		return nil, fmt.Errorf("failed to get room rules: %w", err)
	}

	return rules, nil
}

// GetCleaningSchedule retrieves cleaning schedule for a room
func (r *roomRepository) GetCleaningSchedule(ctx context.Context, roomID uint64) ([]*entity.CleaningSchedule, error) {
	var schedules []*entity.CleaningSchedule
	if err := r.db.WithContext(ctx).
		Where("(room_id = ? OR room_id IS NULL) AND is_active = ?", roomID, true).
		Order("day_of_week, start_time").
		Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("failed to get cleaning schedule: %w", err)
	}

	return schedules, nil
}

// GetRoommates retrieves roommates for a user
func (r *roomRepository) GetRoommates(ctx context.Context, userID uint64) ([]*entity.User, error) {
	// Get user's room
	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.RoomID == nil {
		return []*entity.User{}, nil // Return empty slice if no room assigned
	}

	var roommates []*entity.User
	if err := r.db.WithContext(ctx).
		Where("room_id = ? AND id != ?", *user.RoomID, userID).
		Find(&roommates).Error; err != nil {
		return nil, fmt.Errorf("failed to get roommates: %w", err)
	}

	return roommates, nil
}
