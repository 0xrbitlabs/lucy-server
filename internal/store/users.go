package store

import (
	"fmt"
	"server/internal/types"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
  fmt.Println("came here")
	_, err := u.db.NamedExec(
		`
      insert into users (
        id, user_type, phone_number, password,
        name, profile_picture,
        description, country, town
      )
      values (
        :id, :user_type, :phone_number, :password,
        :name, :profile_picture,
        :description, :country, :town
      )
    `,
		user,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return types.ErrUniqueViolation
			}
		}
		return fmt.Errorf("Failed to insert new user: %w", err)
	}
	return nil
}

func (u *Users) CountByPhoneNumber(phoneNumber string) (int, error) {
	count := 0
	err := u.db.QueryRowx(
		"select count(*) from users where phone_number=$1",
		phoneNumber,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Error while counting phone numbers: %w", err)
	}
	return count, nil
}
