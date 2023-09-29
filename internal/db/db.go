package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var connection *sql.DB

func Connect(dsn string) error {
	if conn, connErr := sql.Open("pgx", dsn); connErr != nil {
		return connErr
	} else {
		connection = conn
	}

	return nil
}

func Connection() *sql.DB {
	return connection
}
