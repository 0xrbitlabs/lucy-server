package repo

import (
	"database/sql"
	"errors"
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

func (r *SessionRepo) GetSessionByID(id string) (*models.Session, error) {
	session := &models.Session{}
	const query = "select * from sessions where id=$1 and valid=true"
	err := r.db.Get(session, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting session by ID: %w", err)
	}
	return session, nil
}
