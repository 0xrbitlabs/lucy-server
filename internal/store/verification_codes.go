package store

import (
	"database/sql"
	"fmt"
	"server/internal/types"

	"github.com/jmoiron/sqlx"
)

type VerificationCodes struct {
	db *sqlx.DB
}

func NewVerificationCodes(db *sqlx.DB) *VerificationCodes {
	return &VerificationCodes{
		db: db,
	}
}

func (vc *VerificationCodes) Insert(code *types.VerificationCode) error {
	_, err := vc.db.NamedExec(
		`
      insert into verification_codes(
        id,
        code,
        sent_to,
        sent_at,
        used
      )
      values (
        :id,
        :code,
        :sent_to,
        :sent_at,
        :used
      )
    `,
		code,
	)
	if err != nil {
		return fmt.Errorf("Error while insering new verification code: %w", err)
	}
	return nil
}

func (vc *VerificationCodes) Get(sent_to, code string) (*types.VerificationCode, error) {
	dbCode := new(types.VerificationCode)
	err := vc.db.Get(
		dbCode,
		"select * from verification_codes where sent_to=$1 and code=$2 and used=false LIMIT 1",
		sent_to,
		code,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrCodeNotFound
		}
		return nil, fmt.Errorf("Error while getting verification code: %w", err)
	}
	return dbCode, nil
}

func (vc *VerificationCodes) SetToUsed(id string) error {
	_, err := vc.db.Exec(
		"update verification_codes set used=true where id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("Error while updating verification code: %w", err)
	}
	return err
}
