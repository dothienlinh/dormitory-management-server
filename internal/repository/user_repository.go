package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Preload("Room.RoomCategory").Preload("EmergencyContacts").First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d not found", user.ID)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Where(user).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with email %s not found", user.Email)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, filter *entity.UserFilter) ([]entity.User, int64, error) {
	filter.Parse()
	whereClause, values := filter.Build()

	var users []entity.User
	var total int64

	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where(whereClause, values...).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).
		Where(whereClause, values...).
		Preload("Room.RoomCategory").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Delete(&entity.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepository) AddUserToRoom(ctx context.Context, payload entity.AddUserToRoom) error {
	var room entity.Room
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Preload("RoomCategory").Where("id = ?", payload.RoomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("room with ID %d not found", payload.RoomID)
		}
		return fmt.Errorf("failed to get room: %w", err)
	}

	var countStudentsInRoom int64
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("room_id = ?", payload.RoomID).Count(&countStudentsInRoom).Error; err != nil {
		return fmt.Errorf("failed to count students in room: %w", err)
	}

	if countStudentsInRoom >= int64(room.RoomCategory.Capacity) {
		return errors.New("room is full")
	}

	var user entity.User
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("id = ? AND role = ?", payload.UserID, entity.UserRoleStudent).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student with ID %d not found", payload.UserID)
		}
		return fmt.Errorf("failed to get student: %w", err)
	}

	if user.RoomID != nil {
		return errors.New("student already has a room")
	}

	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("id = ?", payload.UserID).Update("room_id", room.ID).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) UserLeavesRoom(ctx context.Context, payload entity.UserLeavesRoom) error {
	var room entity.Room
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Where("id = ?", payload.RoomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("room with ID %d not found", payload.RoomID)
		}
		return fmt.Errorf("failed to get room: %w", err)
	}

	var user entity.User
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("id = ? AND role = ?", payload.UserID, entity.UserRoleStudent).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student with ID %d not found", payload.UserID)
		}
		return fmt.Errorf("failed to get student: %w", err)
	}

	if user.RoomID == nil {
		return fmt.Errorf("student not in room")
	}

	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("id = ?", payload.UserID).Update("room_id", nil).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) CheckStudentCodeExists(ctx context.Context, studentCode string) error {
	return r.db.WithContext(ctx).Table(entity.User{}.TableName()).
		Where("student_code = ?", studentCode).
		First(&entity.User{}).Error
}

func (r *userRepository) UpdateUserStatusAccount(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Where("id = ?", user.ID).Update("status_account", user.StatusAccount).Error; err != nil {
		return fmt.Errorf("failed to update user status account: %w", err)
	}
	return nil
}

func (r *userRepository) GetUserByRoles(ctx context.Context, user *entity.User, roles []entity.UserRole) error {
	return r.db.WithContext(ctx).Table(user.TableName()).Where("email = ? AND role IN ?", user.Email, roles).First(user).Error
}

func (r *userRepository) UpdateMe(ctx context.Context, user *entity.User, payload entity.UserUpdateMe) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		updates := make(map[string]interface{})

		if payload.FullName != nil {
			updates["full_name"] = *payload.FullName
		}

		if payload.Phone != nil {
			updates["phone"] = *payload.Phone
		}

		if payload.Birthday != nil {
			updates["birthday"] = *payload.Birthday
		}

		if payload.Gender != "" {
			updates["gender"] = payload.Gender
		}

		if len(updates) > 0 {
			if err := tx.WithContext(ctx).Table(user.TableName()).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update user: %w", err)
			}
		}

		if payload.EmergencyContact.Name != "" && payload.EmergencyContact.Phone != "" {
			var existingContacts []entity.EmergencyContact
			if err := tx.WithContext(ctx).Table(entity.EmergencyContact{}.TableName()).Where("user_id = ?", user.ID).Find(&existingContacts).Error; err != nil {
				return fmt.Errorf("failed to get existing emergency contacts: %w", err)
			}

			if len(existingContacts) > 0 {
				emergencyContactUpdates := map[string]interface{}{
					"name":  payload.EmergencyContact.Name,
					"phone": payload.EmergencyContact.Phone,
				}
				if err := tx.WithContext(ctx).Table(entity.EmergencyContact{}.TableName()).
					Where("user_id = ?", user.ID).
					Updates(emergencyContactUpdates).Error; err != nil {
					return fmt.Errorf("failed to update emergency contact: %w", err)
				}
			} else {
				emergencyContact := &entity.EmergencyContact{
					UserID: user.ID,
					Name:   payload.EmergencyContact.Name,
					Phone:  payload.EmergencyContact.Phone,
				}
				if err := tx.WithContext(ctx).Table(entity.EmergencyContact{}.TableName()).Create(emergencyContact).Error; err != nil {
					return fmt.Errorf("failed to create emergency contact: %w", err)
				}
			}
		}

		return nil
	})
}
