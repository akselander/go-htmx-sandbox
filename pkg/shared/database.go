package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Connect(ConnectionUri string) *sql.DB {
	db, err := sql.Open("sqlite3", ConnectionUri)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	return db
}
