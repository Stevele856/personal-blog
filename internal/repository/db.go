package repository

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

// DB connection
func NewDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	db.SetMaxOpenConns(1)
	return db, nil
}