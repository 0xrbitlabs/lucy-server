package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func PostgresPool() *sqlx.DB {
	pool, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	return pool
}
