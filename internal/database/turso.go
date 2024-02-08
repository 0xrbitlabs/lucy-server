package database

import (
	"database/sql"
	"fmt"
	_ "github.com/libsql/go-libsql"
	"os"
)

func ConnectToTurso() *sql.DB {
	pool, err := sql.Open("libsql", os.Getenv("TURSO_URL"))
	if err != nil {
		panic(err)
	}
	err = pool.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Turso")
	return pool
}
