package store

import (
	"database/sql"
	"fmt"
	"lucy/app_errors"
	"lucy/types"
)

type GetAdminFilter struct {
	Column string
	Value  string
}

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

func (s *Store) GetAdmin(filter GetAdminFilter) (*types.Admin, error) {
	admin := &types.Admin{}
	query := fmt.Sprintf("select * from admins where %s=$1", filter.Column)
	err := s.db.Get(
		admin,
		query,
		filter.Value,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app_errors.ErrResourceNotFound
		}
		return nil, fmt.Errorf("Error while retrieving admin by %s: %w", filter.Column, err)
	}
	return admin, nil
}
