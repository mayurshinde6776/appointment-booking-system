package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB opens a GORM/PostgreSQL connection pool using credentials read from
// the following environment variables:
//
//	DB_HOST     – database host             (default: localhost)
//	DB_PORT     – database port             (default: 5432)
//	DB_USER     – database user             (default: postgres)
//	DB_PASSWORD – database password         (default: postgres)
//	DB_NAME     – database name             (default: appointment_db)
//	DB_SSLMODE  – SSL mode                  (default: disable)
//
// It configures the connection pool, verifies connectivity with a ping, and
// returns the *gorm.DB instance ready for use.
func NewDB() (*gorm.DB, error) {
	dsn := buildDSN()

	// Structured GORM logger that highlights slow queries.
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  resolveLogLevel(),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return nil, fmt.Errorf("database: open connection: %w", err)
	}

	// Tune the underlying *sql.DB connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: retrieve sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	// Verify the server is reachable before the application starts.
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database: ping failed – check DB_* env vars: %w", err)
	}

	log.Printf("database: connected to %s:%s/%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "appointment_db"),
	)

	return db, nil
}

// resolveLogLevel returns Silent in release mode to avoid leaking query data.
func resolveLogLevel() logger.LogLevel {
	if os.Getenv("GIN_MODE") == "release" {
		return logger.Silent
	}
	return logger.Info
}
