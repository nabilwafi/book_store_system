package database

import (
	"fmt"
	"log"
	"os"

	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger.Info("Connected to PostgreSQL database")
}

// Migrate runs auto migration for the database models
func Migrate() {
	err := DB.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Book{},
		&entity.Order{},
		&entity.OrderItem{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	logger.Info("Database migration completed")
}