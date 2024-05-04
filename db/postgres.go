package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func PostgresDB(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
