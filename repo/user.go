package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"lucy/models"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Get(user, "select * from users where phone_number=$1", phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting user by phone number: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Get(user, "select * from users where id=$1")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting user by ID: %w", err)
	}
	return user, nil
}
