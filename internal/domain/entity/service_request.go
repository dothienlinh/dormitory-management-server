package entity

import (
	"time"

	"gorm.io/gorm"
)

type ServiceRequestStatus string

const (
	ServiceRequestStatusPending    ServiceRequestStatus = "pending"
	ServiceRequestStatusApproved   ServiceRequestStatus = "approved"
	ServiceRequestStatusInProgress ServiceRequestStatus = "in_progress"
	ServiceRequestStatusCompleted  ServiceRequestStatus = "completed"
	ServiceRequestStatusRejected   ServiceRequestStatus = "rejected"
)

type ServiceRequestPriority string

const (
	ServiceRequestPriorityLow    ServiceRequestPriority = "low"
	ServiceRequestPriorityMedium ServiceRequestPriority = "medium"
	ServiceRequestPriorityHigh   ServiceRequestPriority = "high"
)

type ServiceRequest struct {
	ID             uint64                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          string                 `json:"title" gorm:"not null"`
	Description    string                 `json:"description" gorm:"type:text"`
	Category       string                 `json:"category" gorm:"not null"`
	Priority       ServiceRequestPriority `json:"priority" gorm:"default:'medium'"`
	Status         ServiceRequestStatus   `json:"status" gorm:"default:'pending'"`
	UserID         uint64                 `json:"user_id" gorm:"not null"`
	User           *User                  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AssignedTo     *uint64                `json:"assigned_to,omitempty"`
	AssignedUser   *User                  `json:"assigned_user,omitempty" gorm:"foreignKey:AssignedTo"`
	CompletionDate *time.Time             `json:"completion_date,omitempty"`
	CreatedAt      time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt         `json:"-" gorm:"index"`
}

func (ServiceRequest) TableName() string {
	return "service_requests"
}
