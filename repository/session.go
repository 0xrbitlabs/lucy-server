package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/lucy/domain"
)

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{db}
}

func (r *SessionRepo) CreateSession(session *domain.Session) error {
	const query = `
    insert into sessions (
      id, valid, user
    )
    values (
      :id, :valid, :user
    )
  `
	_, err := r.db.NamedExec(query, session)
	if err != nil {
		return fmt.Errorf("Error while creating new session")
	}
	return nil
}

func (r *SessionRepo) GetSessionByID(id string) (*domain.Session, error) {
	const query = "select * from sessions where id=$1"
	session := &domain.Session{}
	err := r.db.Get(session, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Error while getting session by ID: %w", err)
	}
	return session, nil
}
