package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const connStr string = "user=test dbname=test sslmode=disable password=1234"

func Connect() error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		username TEXT NOT NULL,
		email TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}
	return nil
}
