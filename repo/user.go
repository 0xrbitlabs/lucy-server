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

func (r *UserRepo) Insert(user *models.User) error {
	const query = `
    insert into users(
      id, username, phone_number,
      password, created_at, account_type
    )
    values (
      :id, :username, :phone_number,
      :password, :created_at, :account_type
    )
  `
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	user := &models.User{}
	const query = "select * from users where phone_number=$1"
	err := r.db.Get(user, query, phoneNumber)
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
	const query = "select * from users where id=$1"
	err := r.db.Get(user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting user by ID: %w", err)
	}
	return user, nil
}

func (r *UserRepo) SetToVerified(id string) error {
	const query = "update users set verified=true where id=$1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error while setting user to verified: %w", err)
	}
	return nil
}
