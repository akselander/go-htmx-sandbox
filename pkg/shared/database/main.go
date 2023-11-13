package database

import (
	"database/sql"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func Connect(ConnectionUri string) (*sql.DB, error) {
	db, err := sql.Open("libsql", ConnectionUri)
	if err != nil {
		return nil, err
	}

	return db, nil
}
