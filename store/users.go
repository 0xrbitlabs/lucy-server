package store

import (
	"database/sql"
	"fmt"
	"lucy/app_errors"
	"lucy/types"
)

func (s *Store) InsertUser(u *types.User) error {
	_, err := s.db.NamedExec(
		`
      insert into users(
        id, type, phone_number, password, username,
        description, country, town
      )
      values(
        :id, :type, :phone_number, :password, :username,
        :description, :country, :town
      )
    `,
		u,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}

func (s *Store) GetUserByID(id string) (*types.User, error) {
	user := &types.User{}
	err := s.db.Get(
		user,
		"select * from users where id==$1",
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app_errors.ErrResourceNotFound
		}
		return nil, fmt.Errorf("Error while retrieving user by id: %w", err)
	}
	return user, nil
}

func (s *Store) CheckForDuplicate(phone string) (int, error) {
	count := 0
	err := s.db.QueryRow(
		`
      select count(*) from users where phone_number=$1
    `,
		phone,
	).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("Error while checking for duplicates: %w", err)
	}
	return count, nil
}
