package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default connection string for local development
		dsn = "host=localhost user=iumicert password=iumicert_secret dbname=iumicert port=5432 sslmode=disable"
	}

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)
	if os.Getenv("ENV") == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		// Optimize for bulk inserts
		CreateBatchSize: 100,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	maxIdleConns := getEnvInt("DB_MAX_IDLE", 10)
	maxOpenConns := getEnvInt("DB_MAX_CONNECTIONS", 100)

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection established")

	// Store in global variable for easy access
	DB = db

	return db, nil
}

// RunMigrations creates/updates database schema
func RunMigrations(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	err := db.AutoMigrate(
		&Student{},
		&Term{},
		&TermReceipt{},
		&AccumulatedReceipt{},
		&VerificationLog{},
		&BlockchainTransaction{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Database migrations completed")
	return nil
}

// CreateIndexes creates additional indexes for performance
func CreateIndexes(db *gorm.DB) error {
	log.Println("🔄 Creating additional indexes...")

	indexes := []string{
		// Composite indexes for common queries
		"CREATE INDEX IF NOT EXISTS idx_term_receipts_student_term ON term_receipts (student_id, term_id)",
		"CREATE INDEX IF NOT EXISTS idx_term_receipts_generated ON term_receipts (generated_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_accumulated_student_type ON accumulated_receipts (student_id, type)",

		// GIN index for JSON queries (if needed)
		"CREATE INDEX IF NOT EXISTS idx_term_receipts_courses ON term_receipts USING gin (revealed_courses)",

		// Verification logs indexes
		"CREATE INDEX IF NOT EXISTS idx_verification_logs_receipt ON verification_logs (receipt_id, verified_at DESC)",
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			log.Printf("⚠️  Warning: Failed to create index: %v", err)
			// Continue with other indexes even if one fails
		}
	}

	log.Println("✅ Indexes created")
	return nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck checks if the database connection is alive
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Helper function to get environment variable as int with default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
