package database

import (
	"appointment-booking/internal/models"

	"gorm.io/gorm"
)

// Migrate runs auto-migration for all domain models.
// It creates or updates tables to match the current model definitions.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Appointment{},
	)
}
