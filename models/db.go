package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected")
}


// pqArray converts a Go slice into a Postgres array type for INSERT/UPDATE
func PqArray(arr []string) interface{} {
	return pq.Array(arr)
}

// pqArrayScan returns a pointer that can scan Postgres arrays into Go slices
func PqArrayScan(dest *[]string) interface{} {
	return pq.Array(dest)
}
