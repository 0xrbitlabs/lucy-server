package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetByPhone(phone string) (*User, error) {
	user := &User{}
	const query = "select * from users where phone=$1"
	err := r.db.Get(user, query, phone)
	if err != nil {
		return nil, fmt.Errorf("Error while getting user by phone: %w", err)
	}
	return user, nil
}

func (r *UserRepo) Insert(user *User) error {
	const query = `
    insert into users(
      id, phone, username, password, account_type
    )
    values(
      :id, :phone, :username, :password, :account_type
    )
  `
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}
