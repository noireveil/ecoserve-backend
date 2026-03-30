package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/noireveil/ecoserve-backend/internal/domain"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database: \n", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS postgis;")

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Technician{},
		&domain.Order{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema: \n", err)
	}

	log.Println("Successfully connected to the PostgreSQL database and migrated schema")
	DB = db
}
