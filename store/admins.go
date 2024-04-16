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
      insert into admins(id, username, password, is_super)
      values(:id, :username, :password, :is_super)
    `,
		admin,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting admin: %w", err)
	}
	return nil
}

func (s *Store) CountAdminByUsername(username string) (int, error) {
	count := 0
	err := s.db.QueryRowx(
		"select count(*) from admins where username=$1",
		username,
	).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("Error while counting admins: %w", err)
	}
	return count, nil
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

func (s *Store) UpdateAdminInfo(username, passwordHash string) error {
	_, err := s.db.Exec("update admins set username=$1, password=$2", username, passwordHash)
	if err != nil {
		return fmt.Errorf("Error while updating admin info: %w", err)
	}
	return nil
}
