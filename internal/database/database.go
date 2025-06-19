package database

import (
	"fmt"
	"log"

	"ishare-task-api/internal/config"
	"ishare-task-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init initializes the database connection and runs migrations
func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connection established and migrations completed")
	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	// Auto migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.AuthorizationCode{},
		&models.AccessToken{},
	)
	if err != nil {
		return err
	}

	// Create indexes for better performance
	if err := createIndexes(db); err != nil {
		return err
	}

	return nil
}

// createIndexes creates database indexes for better performance
func createIndexes(db *gorm.DB) error {
	// User indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return err
	}

	// Task indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at)").Error; err != nil {
		return err
	}

	// Authorization code indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_auth_codes_code ON authorization_codes(code)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_auth_codes_expires_at ON authorization_codes(expires_at)").Error; err != nil {
		return err
	}

	// Access token indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_access_tokens_token ON access_tokens(token)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_access_tokens_expires_at ON access_tokens(expires_at)").Error; err != nil {
		return err
	}

	return nil
} 