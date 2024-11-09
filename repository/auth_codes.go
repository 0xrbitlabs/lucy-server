package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/lucy/domain"
)

type AuthCodeRepo struct {
	db *sqlx.DB
}

func NewAuthCodeRepo(db *sqlx.DB) *AuthCodeRepo {
	return &AuthCodeRepo{db}
}

func (r *AuthCodeRepo) Insert(authCode *domain.AuthCode) error {
	const query = `
    insert into auth_codes(
      code, used, generated_for
    )
    values(
      :code, :used, :generated_for
    )
  `
	_, err := r.db.NamedExec(query, authCode)
	if err != nil {
		return fmt.Errorf("Error while inserting auth code: %w", err)
	}
	return nil
}

func (r *AuthCodeRepo) Get(code, generatedFor string) (*domain.AuthCode, error) {
	const query = `select * from auth_codes where code=$1 and generated_for=$2`
	authCode := &domain.AuthCode{}
	err := r.db.Get(authCode, query, code, generatedFor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Error while getting auth code: %w", err)
	}
	return authCode, nil
}

func (r *AuthCodeRepo) SetToUsed(id int) error {
	_, err := r.db.Exec("update auth_codes set used=true where id=$1", id)
	if err != nil {
		return fmt.Errorf("Error while setting code to used: %w", err)
	}
	return nil
}
