package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetUpDatabaseConnection(logger *slog.Logger) *gorm.DB {
	if err := godotenv.Load(); err != nil {
		if logger != nil {
			logger.Warn("could not load .env file", "err", err)
		}
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbSSL := os.Getenv("DB_SSLMODE")

	// Формируем DSN
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSSL,
	)

	// Подключение через GORM
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize database, got error: %v", err)
	}

	log.Println("✅ Successfully connected to the database")
	return db
}
