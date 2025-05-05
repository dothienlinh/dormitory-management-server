package models

import (
	"dormitory_management/internal/helpers"
	"time"

	"gorm.io/gorm"
)

type IUser interface {
	GenerateStudentCode()
}

type User struct {
	BaseModel
	FullName    string     `json:"full_name" gorm:"type:varchar(255)"`
	StudentCode string     `json:"student_code" gorm:"type:varchar(255);unique"`
	Email       string     `json:"email" gorm:"type:varchar(255);unique"`
	Password    string     `json:"-" gorm:"type:text"`
	Role        UserRole   `json:"role" gorm:"type:user_role;default:student"`
	Gender      UserGender `json:"gender" gorm:"type:user_gender;default:other"`
	Status      UserStatus `json:"status" gorm:"type:user_status;default:active"`
	Phone       *string    `json:"phone" gorm:"type:varchar(255);default:null"`
	Address     *string    `json:"address" gorm:"type:text;default:null"`
	Birthday    *time.Time `json:"birthday" gorm:"default:null"`
	Avatar      *string    `json:"avatar" gorm:"type:text;default:null"`
	RoomID      *uint      `json:"room_id" gorm:"default:null"`
	Room        *Room      `json:"room" gorm:"foreignKey:RoomID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	u.StudentCode = time.Now().Format("20060102150405")

	return nil
}

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
)

type UserRegister struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
