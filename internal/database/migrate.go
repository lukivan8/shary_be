package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// RunMigrations runs database migrations using golang-migrate
func RunMigrations(db *sqlx.DB, migrationsPath string) error {
	// Get the underlying sql.DB from sqlx
	sqlDB := db.DB

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(db *sqlx.DB, migrationsPath string) (uint, error) {
	sqlDB := db.DB

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return 0, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		return 0, fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		return 0, fmt.Errorf("database is in dirty state at version %d", version)
	}

	return version, nil
}
