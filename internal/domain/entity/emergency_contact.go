package entity

type EmergencyContact struct {
	Base
	UserID       uint64 `json:"user_id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Relationship string `json:"relationship"`
}

func (EmergencyContact) TableName() string {
	return "emergency_contacts"
}
