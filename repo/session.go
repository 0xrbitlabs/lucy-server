package repo

import (
	"fmt"
	"lucy/models"

	"github.com/jmoiron/sqlx"
)

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{db}
}

func (r *SessionRepo) Insert(session *models.Session) error {
	const query = `
    INSERT INTO SESSIONS (
      id, user_id
    )
    VALUES(
      :id, :user_id
    )
  `
	_, err := r.db.NamedExec(query, session)
	if err != nil {
		return fmt.Errorf("Error while inserting session: %w", err)
	}
	return nil
}
