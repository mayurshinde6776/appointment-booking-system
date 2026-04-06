package models

import (
	"time"

	"gorm.io/gorm"
)

// AppointmentStatus represents the lifecycle state of an appointment.
type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "pending"
	StatusConfirmed AppointmentStatus = "confirmed"
	StatusCancelled AppointmentStatus = "cancelled"
	StatusCompleted AppointmentStatus = "completed"
)

// Appointment is the core domain model persisted in the database.
type Appointment struct {
	ID          uint              `gorm:"primaryKey;autoIncrement"                   json:"id"`
	PatientName string            `gorm:"type:varchar(255);not null"                 json:"patient_name"   validate:"required,min=2,max=255"`
	DoctorName  string            `gorm:"type:varchar(255);not null"                 json:"doctor_name"    validate:"required,min=2,max=255"`
	Date        time.Time         `gorm:"not null"                                   json:"date"           validate:"required"`
	Duration    int               `gorm:"not null;default:30"                        json:"duration_mins"  validate:"required,min=10,max=480"`
	Status      AppointmentStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Notes       string            `gorm:"type:text"                                  json:"notes"`
	CreatedAt   time.Time         `                                                  json:"created_at"`
	UpdatedAt   time.Time         `                                                  json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index"                                      json:"-"`
}
