package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

var db *sqlx.DB

func Initialize() error {
	var envFile string

	// Check if running in production environment
	if os.Getenv("ENV") == "production" {
		envFile = "/app/.env" // Production .env path
	} else {
		envFile = "../../.env" // Default to this path for local development
	}

	err := godotenv.Load(envFile)
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connectionStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=%s", dbUser, dbName, dbPassword, dbHost, dbSSLMode)

	db, err = sqlx.Connect("postgres", connectionStr)
	if err != nil {
		return err
	}
	return db.Ping()
}

func GetDB() *sqlx.DB {
	return db
}
