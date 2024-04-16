package store

import (
	"database/sql"
	"fmt"
	"lucy/app_errors"
	"lucy/types"

	"github.com/lib/pq"
)

func (s *Store) InsertCategory(category types.Category) error {
	_, err := s.db.NamedExec(
		`
      insert into categories(id, label, description, enabled)
      values(:id, :label, :description, :enabled)
    `,
		category,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return app_errors.ErrDuplicateResource
			}
		}
		return fmt.Errorf("Error while inserting category: %w", err)
	}
	return nil
}

// Returns all categories if called by admin
func (s *Store) GetCategories(getAll bool) ([]types.Category, error) {
	data := []types.Category{}
	query := "select * from categories where enabled=true"
	if getAll {
		query = "select * from categories"
	}
	err := s.db.Select(
		&data,
		query,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return data, nil
		}
		return data, fmt.Errorf("Error while retrieving categories: %w", err)
	}
	return data, nil
}
