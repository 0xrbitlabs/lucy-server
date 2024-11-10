package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/lucy/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetByPhone(phone string) (*domain.User, error) {
	user := &domain.User{}
	const query = "select * from users where phone=$1"
	err := r.db.Get(user, query, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Error while getting user by phone: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetAll() ([]domain.User, error) {
	data := make([]domain.User, 0)
	const query = "select * from users"
	err := r.db.Select(&data, query)
	if err != nil {
		return nil, fmt.Errorf("Error while getting all users: %w", err)
	}
	return data, nil
}

func (r *UserRepo) Insert(tx *sqlx.Tx, user *domain.User) error {
	const query = `
    insert into users(
      id, phone, username, password, account_type
    )
    values(
      :id, :phone, :username, :password, :account_type
    )
  `
	var err error
	if tx != nil {
		_, err = tx.NamedExec(query, user)
	} else {
		_, err = r.db.NamedExec(query, user)
	}
	if err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateUser(user *domain.User) error {
	const query = `
    update users set username=$1, password=$2
    where id=$3
  `
	_, err := r.db.Exec(query, user.Username, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("Error while updating user data: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateUserPassword(userID, password string) error {
	const query = `
    update users set password=$1
    where id=$2
  `
	_, err := r.db.Exec(query, password, userID)
	if err != nil {
		return fmt.Errorf("Error while updating user password: %w", err)
	}
	return nil
}
