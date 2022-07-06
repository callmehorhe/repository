package repository

import (
	"database/sql"
	"fmt"
)

// connString returns connection string for db.
func connString(host, port, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// NewPostgresDB creates new connection to DB.
func NewPostgresDB() (*sql.DB, error) {
	conn := connString("localhost", "5432", "postgres", "12345", "payments")
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
