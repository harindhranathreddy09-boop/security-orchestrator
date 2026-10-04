package db

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
)

// NewPostgres returns a *sql.DB ready to use.
// The caller should defer db.Close().
func NewPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// Basic connection test
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	// Reasonable defaults for a short‑lived service
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(30 * time.Minute)
	return db, nil
}
