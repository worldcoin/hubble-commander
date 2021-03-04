package db

import (
	"github.com/jmoiron/sqlx"
)

type TransactionController struct {
	tx       *sqlx.Tx
	isLocked bool
}

func (t *TransactionController) Rollback() error {
	if !t.isLocked {
		t.isLocked = true
		return t.tx.Rollback()
	}
	return nil
}

func (t *TransactionController) Commit() error {
	if !t.isLocked {
		t.isLocked = true
		return t.tx.Commit()
	}
	return nil
}
