package database

import (
	"os"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresPool() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("PG_URL"))
	if err != nil {
		panic(err)
	}
	return db
}
