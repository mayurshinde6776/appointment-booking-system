package models

import "time"

// User represents a patient or end-user who can book appointments.
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	Name      string    `gorm:"type:varchar(255);not null"            json:"name"       validate:"required,min=2,max=255"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"      validate:"required,email"`
	CreatedAt time.Time `gorm:"autoCreateTime"                        json:"created_at"`

	// Relationships
	Bookings []Booking `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bookings,omitempty"`
}
