package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"lucy/models"

	"github.com/jmoiron/sqlx"
)

type VerificationCodeRepo struct {
	db *sqlx.DB
}

func NewVerificationCodeRepo(db *sqlx.DB) *VerificationCodeRepo {
	return &VerificationCodeRepo{db}
}

func (r *VerificationCodeRepo) Insert(vc *models.VerificationCode) error {
	const query = `
    insert into verification_codes (
      code, generated_for, generated_at
    )
    values (
      :code, :generated_for, :generated_at
    )
  `
	_, err := r.db.NamedExec(query, vc)
	if err != nil {
		return fmt.Errorf("Error while inserting verification code: %w", err)
	}
	return nil
}

func (r *VerificationCodeRepo) SetToUsed(id int) error {
	const query = "update verification_codes set used=true where id=$1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error while setting verification code to used: %w", err)
	}
	return nil
}

func (r *VerificationCodeRepo) GetByCode(code, userID string) (*models.VerificationCode, error) {
	vc := &models.VerificationCode{}
	const query = `
    select * from verification_codes where code=$1 and generated_for=$2 and used=false
  `
	err := r.db.Get(vc, query, code, userID)
	if err != nil {
    if errors.Is(err, sql.ErrNoRows){
      return nil, nil
    }
		return nil, fmt.Errorf("Error while getting verification code: %w", err)
	}
	return vc, nil
}
