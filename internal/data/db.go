package data

import (
	"fmt"
	"log"
	"time"

	"toko-emas/internal/conf"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a new database connection
func NewDB(c *conf.Config) (*gorm.DB, error) {
	dsn := c.Data.Database.DSN
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is empty. Please set DATABASE_URL or database configuration")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Recommended for cloud proxies/pgx (like Railway)
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Optimize connection pool for cloud and production environments
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return db, nil
}
