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
      insert into users (
        id, user_type, phone_number, password,
        verified, name, profile_picture,
        description, country, town
      )
      values (
        :id, :user_type, :phone_number, :password,
        :verified, :name, :profile_picture,
        :description, :country, :town
      )
    `,
		user,
	)
	if err != nil {
		return fmt.Errorf("Failed to insert new user: %w", err)
	}
	return nil
}
