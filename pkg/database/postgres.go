package database

import (
	"fmt"
	"os"
)

// buildDSN constructs a PostgreSQL DSN string from environment variables.
// Called by NewDB in db.go.
func buildDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "appointment_db"),
		getEnv("DB_SSLMODE", "disable"),
	)
}

// getEnv returns the value of the environment variable named by key.
// If the variable is unset or empty, fallback is returned.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
