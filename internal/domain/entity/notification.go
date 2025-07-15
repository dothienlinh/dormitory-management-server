package entity

import (
	"time"

	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypePayment      NotificationType = "payment"
	NotificationTypeMaintenance  NotificationType = "maintenance"
	NotificationTypeEvent        NotificationType = "event"
	NotificationTypeService      NotificationType = "service"
	NotificationTypeAnnouncement NotificationType = "announcement"
)

type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityMedium NotificationPriority = "medium"
	NotificationPriorityHigh   NotificationPriority = "high"
)

type Notification struct {
	ID        uint64               `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string               `json:"title" gorm:"not null"`
	Content   string               `json:"content" gorm:"type:text"`
	Type      NotificationType     `json:"type" gorm:"not null"`
	Priority  NotificationPriority `json:"priority" gorm:"default:'medium'"`
	IsRead    bool                 `json:"is_read" gorm:"default:false"`
	UserID    uint64               `json:"user_id" gorm:"not null"`
	User      *User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ReadAt    *time.Time           `json:"read_at,omitempty"`
	CreatedAt time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt       `json:"-" gorm:"index"`
}

func (Notification) TableName() string {
	return "notifications"
}

type RoomIssue struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Title        string         `json:"title" gorm:"not null"`
	Description  string         `json:"description" gorm:"type:text"`
	Category     string         `json:"category" gorm:"not null"`
	Priority     string         `json:"priority" gorm:"default:'medium';check:priority IN ('low','medium','high')"`
	Status       string         `json:"status" gorm:"default:'pending';check:status IN ('pending','in_progress','resolved','rejected')"`
	RoomID       uint64         `json:"room_id" gorm:"not null"`
	Room         *Room          `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	ReportedBy   uint64         `json:"reported_by" gorm:"not null"`
	Reporter     *User          `json:"reporter,omitempty" gorm:"foreignKey:ReportedBy"`
	AssignedTo   *uint64        `json:"assigned_to,omitempty"`
	AssignedUser *User          `json:"assigned_user,omitempty" gorm:"foreignKey:AssignedTo"`
	ResolvedDate *time.Time     `json:"resolved_date,omitempty"`
	Resolution   string         `json:"resolution" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (RoomIssue) TableName() string {
	return "room_issues"
}

type RoomRule struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	Category    string         `json:"category" gorm:"not null;check:category IN ('general','safety','hygiene','behavior')"`
	Priority    int            `json:"priority" gorm:"default:1"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (RoomRule) TableName() string {
	return "room_rules"
}

type CleaningSchedule struct {
	ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	RoomID    *uint64        `json:"room_id,omitempty"` // null means applies to all rooms
	Room      *Room          `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	DayOfWeek string         `json:"day_of_week" gorm:"not null;check:day_of_week IN ('monday','tuesday','wednesday','thursday','friday','saturday','sunday')"`
	StartTime string         `json:"start_time" gorm:"not null"`
	EndTime   string         `json:"end_time" gorm:"not null"`
	Type      string         `json:"type" gorm:"not null;check:type IN ('regular','deep','inspection')"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (CleaningSchedule) TableName() string {
	return "cleaning_schedules"
}

type StudentRoomDetails struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	BuildingName     string              `json:"building_name"`
	BuildingAddress  string              `json:"building_address"`
	Floor            int                 `json:"floor"`
	RoomType         string              `json:"room_type"`
	Capacity         int                 `json:"capacity"`
	Area             float64             `json:"area"`
	MonthlyPrice     float64             `json:"monthly_price"`
	CurrentOccupants int                 `json:"current_occupants"`
	Status           string              `json:"status"`
	Facilities       []*Amenity          `json:"facilities"`
	Residents        []*User             `json:"residents"`
	Rules            []*RoomRule         `json:"rules"`
	RecentIssues     []*RoomIssue        `json:"recent_issues"`
	CleaningSchedule []*CleaningSchedule `json:"cleaning_schedule"`
	BillHistory      []*Payment          `json:"bill_history"`
	Room             *Room               `json:"room,omitempty"`
}

type RoomStats struct {
	TotalFacilities   int    `json:"total_facilities"`
	WorkingFacilities int    `json:"working_facilities"`
	PendingIssues     int    `json:"pending_issues"`
	ResolvedIssues    int    `json:"resolved_issues"`
	CurrentOccupancy  int    `json:"current_occupancy"`
	NextCleaning      string `json:"next_cleaning"`
}
