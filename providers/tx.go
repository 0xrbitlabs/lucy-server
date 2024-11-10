package providers

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TxProvider struct {
	db *sqlx.DB
}

func NewTxProvider(db *sqlx.DB) *TxProvider {
	return &TxProvider{db}
}

func (p *TxProvider) Provide() (*sqlx.Tx, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("Error while getting transaction: %w", err)
	}
	return tx, nil
}
