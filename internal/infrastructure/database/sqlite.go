package database

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/artorias501/blog-service/internal/infrastructure/persistence/model"
)

// Config holds the database configuration parameters
type Config struct {
	// DSN is the data source name (file path for SQLite)
	DSN string

	// LogLevel controls the logging level (Silent, Error, Warn, Info)
	LogLevel logger.LogLevel

	// Connection pool settings
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		DSN:             "blog.db",
		LogLevel:        logger.Warn,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// NewConnection creates a new SQLite database connection with the given configuration
func NewConnection(cfg Config) (*gorm.DB, error) {
	// Open SQLite connection
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// NewConnectionWithMigrate creates a new connection and runs auto-migration
func NewConnectionWithMigrate(cfg Config) (*gorm.DB, error) {
	db, err := NewConnection(cfg)
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&model.PostModel{},
		&model.TagModel{},
		&model.CommentModel{},
		&model.PostTagModel{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}
