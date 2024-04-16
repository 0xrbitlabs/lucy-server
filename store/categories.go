package store

import (
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
