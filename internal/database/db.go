package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	fmt.Println("Connected to PostgreSQL")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
