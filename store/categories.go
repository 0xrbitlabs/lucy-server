package store

import (
	"fmt"
	"lucy/types"
)

func (s *Store) InsertCategory(category types.Category) error {
	_, err := s.db.NamedExec(
		`
      insert into categories(id, label, description)
      values(:id, :label, :description)
    `,
		category,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting category: %w", err)
	}
	return nil
}
