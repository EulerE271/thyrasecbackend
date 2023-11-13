package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

var db *sqlx.DB

func Initialize() error {
	connectionStr := "user=postgres dbname=thyrasec password=root host=localhost sslmode=disable"
	var err error
	db, err = sqlx.Connect("postgres", connectionStr)
	if err != nil {
		return err
	}
	return db.Ping()
}

func GetDB() *sqlx.DB {
	return db
}
