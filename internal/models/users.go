package models

import (
	"dormitory_management/internal/helpers"
	"strings"
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
	Birthday    *time.Time `json:"birthday" gorm:"default:null"`
	Avatar      *string    `json:"avatar" gorm:"type:text;default:null"`
	RoomRentID  *uint      `json:"-" gorm:"default:null"`
	RoomRent    *RoomRent  `json:"room_rent"`
	ContractID  *uint      `json:"-" gorm:"default:null"`
	Contract    *Contract  `json:"contract"`
}

type UserSimple struct {
	BaseModel
	FullName    string     `json:"full_name"`
	StudentCode string     `json:"student_code"`
	Email       string     `json:"email"`
	Gender      UserGender `json:"gender"`
	Status      UserStatus `json:"status"`
	Phone       *string    `json:"phone"`
	Birthday    *time.Time `json:"birthday"`
	Avatar      *string    `json:"avatar"`
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
	UserStatusAbsent   UserStatus = "absent"
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

type FilterUser struct {
	Status  UserStatus `form:"status" validate:"omitempty,oneof=active inactive absent"`
	Keyword string     `form:"keyword"`
	Gender  UserGender `form:"gender" validate:"omitempty,oneof=male female other"`
	Pagination
}

func (f FilterUser) Build() (string, []interface{}) {
	conditions := []string{"role = ?"}
	values := []interface{}{UserRoleStudent}

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
