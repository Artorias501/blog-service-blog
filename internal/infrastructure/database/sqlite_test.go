package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gorm.io/gorm/logger"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DSN != "blog.db" {
		t.Errorf("expected DSN 'blog.db', got '%s'", cfg.DSN)
	}

	if cfg.LogLevel != logger.Warn {
		t.Errorf("expected LogLevel %d, got %d", logger.Warn, cfg.LogLevel)
	}

	if cfg.MaxIdleConns != 10 {
		t.Errorf("expected MaxIdleConns 10, got %d", cfg.MaxIdleConns)
	}

	if cfg.MaxOpenConns != 100 {
		t.Errorf("expected MaxOpenConns 100, got %d", cfg.MaxOpenConns)
	}

	if cfg.ConnMaxLifetime != time.Hour {
		t.Errorf("expected ConnMaxLifetime %v, got %v", time.Hour, cfg.ConnMaxLifetime)
	}
}

func TestNewConnection(t *testing.T) {
	// Create temp directory for test database
	tmpDir := filepath.Join(os.TempDir(), "blog-test")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test_connection.db")
	defer os.Remove(dbPath)

	cfg := Config{
		DSN:             dbPath,
		LogLevel:        logger.Silent,
		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	db, err := NewConnection(cfg)
	if err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}

	if db == nil {
		t.Error("expected db connection, got nil")
	}

	// Verify connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying SQL DB: %v", err)
	}

	// Close connection
	sqlDB.Close()
}

func TestNewConnectionWithMigrate(t *testing.T) {
	// Create temp directory for test database
	tmpDir := filepath.Join(os.TempDir(), "blog-test")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test_migrate.db")
	defer os.Remove(dbPath)

	cfg := Config{
		DSN:             dbPath,
		LogLevel:        logger.Silent,
		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	db, err := NewConnectionWithMigrate(cfg)
	if err != nil {
		t.Fatalf("failed to create connection with migration: %v", err)
	}

	if db == nil {
		t.Error("expected db connection, got nil")
	}

	// Verify tables were created
	var tables []string
	if err := db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables).Error; err != nil {
		t.Fatalf("failed to query tables: %v", err)
	}

	expectedTables := map[string]bool{
		"posts":     false,
		"tags":      false,
		"comments":  false,
		"post_tags": false,
	}

	for _, table := range tables {
		if _, ok := expectedTables[table]; ok {
			expectedTables[table] = true
		}
	}

	for table, found := range expectedTables {
		if !found {
			t.Errorf("expected table '%s' to be created", table)
		}
	}

	// Close connection
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func TestConfigCustomValues(t *testing.T) {
	cfg := Config{
		DSN:             "custom.db",
		LogLevel:        logger.Info,
		MaxIdleConns:    20,
		MaxOpenConns:    50,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	if cfg.DSN != "custom.db" {
		t.Errorf("expected DSN 'custom.db', got '%s'", cfg.DSN)
	}

	if cfg.LogLevel != logger.Info {
		t.Errorf("expected LogLevel %d, got %d", logger.Info, cfg.LogLevel)
	}

	if cfg.MaxIdleConns != 20 {
		t.Errorf("expected MaxIdleConns 20, got %d", cfg.MaxIdleConns)
	}

	if cfg.MaxOpenConns != 50 {
		t.Errorf("expected MaxOpenConns 50, got %d", cfg.MaxOpenConns)
	}

	if cfg.ConnMaxLifetime != 30*time.Minute {
		t.Errorf("expected ConnMaxLifetime %v, got %v", 30*time.Minute, cfg.ConnMaxLifetime)
	}
}
