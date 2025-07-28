package entity

import (
	"dormitory_management/internal/helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleStudent UserRole = "student"
	UserRoleAdmin   UserRole = "admin"
	UserRoleStaff   UserRole = "staff"
)

type UserGender string

const (
	UserGenderMale   UserGender = "male"
	UserGenderFemale UserGender = "female"
	UserGenderOther  UserGender = "other"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusAbsent   UserStatus = "absent"
)

type StatusAccount string

const (
	StatusAccountPending  StatusAccount = "pending"
	StatusAccountApproved StatusAccount = "approved"
	StatusAccountRejected StatusAccount = "rejected"
	StatusAccountBanned   StatusAccount = "banned"
)

type User struct {
	ID                uint64              `json:"id" gorm:"primarykey"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	DeletedAt         gorm.DeletedAt      `json:"-" gorm:"index"`
	FullName          string              `json:"full_name" gorm:"not null"`
	StudentCode       *string             `json:"student_code"`
	Email             string              `json:"email" gorm:"uniqueIndex;not null"`
	Password          string              `json:"-" gorm:"not null"`
	Role              UserRole            `json:"role" gorm:"not null;default:'student'"`
	Gender            UserGender          `json:"gender" gorm:"not null"`
	Status            UserStatus          `json:"status" gorm:"not null;default:'active'"`
	Phone             string              `json:"phone" gorm:"not null"`
	IsVerify          bool                `json:"is_verify" gorm:"default:false"`
	StatusAccount     StatusAccount       `json:"status_account" gorm:"not null;default:'pending'"`
	Birthday          *time.Time          `json:"birthday"`
	Avatar            *string             `json:"avatar"`
	Address           *string             `json:"address"`
	RoomID            *uint               `json:"room_id"`
	Room              *Room               `json:"room" gorm:"foreignKey:RoomID"`
	Contracts         *[]Contract         `json:"contracts" gorm:"foreignKey:UserID"`
	Payments          *[]Payment          `json:"payments" gorm:"foreignKey:UserId"`
	EmergencyContacts *[]EmergencyContact `json:"emergency_contacts" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

type UpdateUser struct {
	FullName    *string     `json:"full_name"`
	StudentCode *string     `json:"student_code"`
	Email       *string     `json:"email"`
	Password    *string     `json:"-"`
	Role        *UserRole   `json:"role"`
	Gender      *UserGender `json:"gender"`
	Status      *UserStatus `json:"status"`
	Phone       *string     `json:"phone"`
	Birthday    *time.Time  `json:"birthday"`
	Avatar      *string     `json:"avatar"`
}

type UserSimple struct {
	ID          uint64         `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	FullName    string         `json:"full_name"`
	StudentCode *string        `json:"student_code"`
	Email       string         `json:"email"`
	Gender      UserGender     `json:"gender"`
	Status      UserStatus     `json:"status"`
	Phone       string         `json:"phone"`
	Birthday    *time.Time     `json:"birthday"`
	Avatar      *string        `json:"avatar"`
}

func (u *User) ToSimple() UserSimple {
	return UserSimple{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
		FullName:    u.FullName,
		StudentCode: u.StudentCode,
		Email:       u.Email,
		Gender:      u.Gender,
		Status:      u.Status,
		Phone:       u.Phone,
		Birthday:    u.Birthday,
		Avatar:      u.Avatar,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := helper.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	return nil
}

func (u *User) ResetPassword(newPassword string) error {
	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

type UserFilter struct {
	Status        UserStatus    `form:"status" binding:"omitempty,oneof=active inactive absent"`
	Keyword       string        `form:"keyword" binding:"omitempty"`
	Gender        UserGender    `form:"gender" binding:"omitempty,oneof=male female other"`
	Role          UserRole      `form:"role" binding:"omitempty,oneof=student staff"`
	StatusAccount StatusAccount `form:"status_account" binding:"omitempty,oneof=pending approved rejected banned"`
	HasRoom       *bool         `form:"has_room" binding:"omitempty"`
	Pagination
}

func (f UserFilter) Build() (string, []interface{}) {
	conditions := []string{"is_verify = ?"}
	values := []interface{}{true}

	if f.Status != "" {
		conditions = append(conditions, "status = ?")
		values = append(values, f.Status)
	}

	if f.Gender != "" {
		conditions = append(conditions, "gender = ?")
		values = append(values, f.Gender)
	}

	if f.Keyword != "" {
		conditions = append(conditions, "(full_name LIKE ? OR email LIKE ? OR phone LIKE ? OR student_code LIKE ?)")
		keyword := "%" + f.Keyword + "%"
		values = append(values, keyword, keyword, keyword, keyword)
	}

	if f.Role != "" {
		conditions = append(conditions, "role = ?")
		values = append(values, f.Role)
	} else {
		conditions = append(conditions, "role != ?")
		values = append(values, UserRoleAdmin)
	}

	if f.HasRoom != nil {
		if !*f.HasRoom {
			conditions = append(conditions, "room_id IS NULL")
		} else if *f.HasRoom {
			conditions = append(conditions, "room_id IS NOT NULL")
		}
	}

	if f.StatusAccount != "" {
		conditions = append(conditions, "status_account = ?")
		values = append(values, f.StatusAccount)
	}

	whereClause := strings.Join(conditions, " AND ")

	return whereClause, values
}

type UserRegister struct {
	FullName string   `json:"full_name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Phone    string   `json:"phone" binding:"required"`
	Password string   `json:"password" binding:"required,min=8"`
	Role     UserRole `json:"role" binding:"required,oneof=student staff"`
}

type LoginType string

const (
	LoginTypeStudent LoginType = "student"
	LoginTypeManager LoginType = "manager"
)

type UserLogin struct {
	Email    string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required,min=8"`
	Type     LoginType `json:"type" binding:"required,oneof=student manager"`
}

type UserRefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserDTO struct {
	ID          uint       `json:"id"`
	FullName    string     `json:"full_name"`
	StudentCode *string    `json:"student_code"`
	Email       string     `json:"email"`
	Role        UserRole   `json:"role"`
	Gender      UserGender `json:"gender"`
	Status      UserStatus `json:"status"`
	Phone       string     `json:"phone"`
	Birthday    *time.Time `json:"birthday"`
	Avatar      *string    `json:"avatar"`
	Room        *Room      `json:"room"`
	Contract    *Contract  `json:"contract"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type AddUserToRoom struct {
	UserID uint64 `json:"user_id" binding:"required,numeric"`
	RoomID uint   `json:"room_id" binding:"required,numeric"`
}

type UserLeavesRoom struct {
	UserID uint64 `json:"user_id" binding:"required,numeric"`
	RoomID uint   `json:"room_id" binding:"required,numeric"`
}

type UserStatusAccountUpdate struct {
	StatusAccount StatusAccount `json:"status_account" binding:"required,oneof=pending approved rejected banned"`
}

type UserUpdateMe struct {
	FullName         *string    `json:"full_name" binding:"required"`
	Phone            *string    `json:"phone" binding:"required"`
	Birthday         *time.Time `json:"birthday" binding:"omitempty"`
	Gender           UserGender `json:"gender" binding:"omitempty,oneof=male female other"`
	Address          *string    `json:"address" binding:"omitempty"`
	EmergencyContact struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone" binding:"required"`
	} `json:"emergency_contact" binding:"omitempty"`
}
