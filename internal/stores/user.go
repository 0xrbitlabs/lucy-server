package stores

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"server/internal/errors"
	"server/internal/models"
)

type Users struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *Users {
	return &Users{
		db: db,
	}
}

func (u *Users) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	user := new(models.User)
	err := u.db.Get(user, "select * from users where phone_number=$1", phoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.ErrUserNotFound
		}
		return user, fmt.Errorf("Error while querying user with phone number: %w", err)
	}
	return user, nil
}
