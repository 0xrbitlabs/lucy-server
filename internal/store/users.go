package store

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"server/internal/types"
)

type Users struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *Users {
	return &Users{
		db: db,
	}
}

func (u *Users) Insert(user *types.User) error {
	_, err := u.db.NamedExec(
		`
      insert into users(id, phone_number, full_name, profile_picture)
      values(:id, :phone_number, :full_name, :profile_picture)
    `,
		user,
	)
	if err != nil {
		return fmt.Errorf("Failed to insert new user: %w", err)
	}
	return nil
}
