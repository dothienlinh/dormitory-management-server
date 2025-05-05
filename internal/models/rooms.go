package models

type Room struct {
	BaseModel
	Capacity    int        `json:"capacity" gorm:"type:int;not null"`
	Price       int        `json:"price" gorm:"type:int;not null"`
	RoomType    RoomType   `json:"room_type" gorm:"type:room_type;not null"`
	Status      RoomStatus `json:"status" gorm:"type:room_status;not null"`
	Acreage     int        `json:"acreage" gorm:"type:int;not null"`
	Description string     `json:"description" gorm:"type:text;not null"`
	Users       []User     `json:"users" gorm:"foreignKey:RoomID"`
}

type RoomType string

const (
	Standard RoomType = "standard"
	Premium  RoomType = "premium"
)

type RoomStatus string

const (
	Available   RoomStatus = "available"
	Occupied    RoomStatus = "occupied"
	Maintenance RoomStatus = "maintenance"
)
