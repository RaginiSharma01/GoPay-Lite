package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	var err error
	dbURL := os.Getenv("DATABASE_URL") // ✅ Correct key

	log.Println("Connecting to DB with URL:", dbURL)

	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

func Close() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
