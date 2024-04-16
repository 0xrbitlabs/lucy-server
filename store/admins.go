package store

import (
	"database/sql"
	"fmt"
	"lucy/app_errors"
	"lucy/types"
)

func (s *Store) InsertAdmin(admin types.Admin) error {
	_, err := s.db.NamedExec(
		`
      insert into admins(username, password, is_super)
      values(:username, :password, :is_super)
    `,
		admin,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting admin: %w", err)
	}
	return nil
}

func (s *Store) GetAdminByUsername(username string) (*types.Admin, error) {
	admin := &types.Admin{}
	err := s.db.Get(
		admin,
		"select * from admins where username=$1",
		username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app_errors.ErrResourceNotFound
		}
		return nil, fmt.Errorf("Error while retrieving admin by username: %w", err)
	}
	return admin, nil
}
