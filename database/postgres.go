package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetPostgresPool() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
  fmt.Println("Connected to Postgres")
	return db
}
