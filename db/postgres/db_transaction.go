package postgres

import (
	"fmt"
)

type rawController interface {
	Rollback() error
	Commit() error
}

type TransactionController struct {
	tx       rawController
	isLocked bool
}

// nolint:gocritic
func (t *TransactionController) Rollback(cause *error) {
	if !t.isLocked {
		t.isLocked = true
		if rollbackErr := t.tx.Rollback(); rollbackErr != nil {
			*cause = fmt.Errorf("rollback caused by: %w, failed with: %v", *cause, rollbackErr)
		}
	}
}

func (t *TransactionController) Commit() error {
	if !t.isLocked {
		t.isLocked = true
		return t.tx.Commit()
	}
	return nil
}
