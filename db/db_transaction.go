package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TransactionController struct {
	tx       *sqlx.Tx
	isLocked bool
}

func (t *TransactionController) Rollback(cause error) error {
	if !t.isLocked {
		t.isLocked = true
		if rollbackErr := t.tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rollback caused by: %w, failed with: %v", cause, rollbackErr)
		}
	}
	return cause
}

func (t *TransactionController) Commit() error {
	if !t.isLocked {
		t.isLocked = true
		return t.tx.Commit()
	}
	return nil
}
