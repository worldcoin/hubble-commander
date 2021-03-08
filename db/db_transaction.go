package db

import (
	"github.com/jmoiron/sqlx"
)

type TransactionController struct {
	tx       *sqlx.Tx
	isLocked bool
}

func (t *TransactionController) Rollback() {
	if !t.isLocked {
		t.isLocked = true
		err := t.tx.Rollback()
		if err != nil {
			panic(err)
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
