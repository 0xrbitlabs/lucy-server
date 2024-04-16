package store

import (
	"lucy/database"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore() *Store {
	db := database.GetPostgresPool()
	return &Store{
		db: db,
	}
}
