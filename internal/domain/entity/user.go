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

type User struct {
	Base
	FullName    string      `json:"full_name"`
	StudentCode string      `json:"student_code"`
	Email       string      `json:"email"`
	Password    string      `json:"-"`
	Role        UserRole    `json:"role"`
	Gender      UserGender  `json:"gender"`
	Status      UserStatus  `json:"status"`
	Phone       string      `json:"phone"`
	IsVerify    bool        `json:"is_verify"`
	Birthday    *time.Time  `json:"birthday"`
	Avatar      *string     `json:"avatar"`
	RoomID      *uint       `json:"-"`
	Room        *Room       `json:"room"`
	Contracts   *[]Contract `json:"contracts"`
	Payments    *[]Payment  `json:"payments"`
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
	Base
	FullName    string     `json:"full_name"`
	StudentCode string     `json:"student_code"`
	Email       string     `json:"email"`
	Gender      UserGender `json:"gender"`
	Status      UserStatus `json:"status"`
	Phone       string     `json:"phone"`
	Birthday    *time.Time `json:"birthday"`
	Avatar      *string    `json:"avatar"`
}

func (u *User) ToSimple() UserSimple {
	return UserSimple{
		Base:        u.Base,
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

	u.StudentCode = time.Now().Format("20060102150405")

	return nil
}

type UserFilter struct {
	Status  UserStatus `form:"status" binding:"omitempty,oneof=active inactive absent"`
	Keyword string     `form:"keyword"`
	Gender  UserGender `form:"gender" binding:"omitempty,oneof=male female other"`
	Pagination
}

func (f UserFilter) Build() (string, []interface{}) {
	conditions := []string{"role = ?", "is_verify = ?"}
	values := []interface{}{UserRoleStudent, true}

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

	whereClause := strings.Join(conditions, " AND ")

	return whereClause, values
}

type UserRegister struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserRefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserDTO struct {
	ID          uint       `json:"id"`
	FullName    string     `json:"full_name"`
	StudentCode string     `json:"student_code"`
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
