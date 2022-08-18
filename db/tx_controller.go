package db

import (
	"fmt"

	"github.com/pkg/errors"
)

type rawController interface {
	Rollback() error
	Commit() error
}

type TxController struct {
	tx       rawController
	isLocked bool
}

func NewTxController(tx rawController, isLocked bool) *TxController {
	return &TxController{tx, isLocked}
}

//nolint:gocritic
func (t *TxController) Rollback(cause *error) {
	if !t.isLocked {
		t.isLocked = true
		if rollbackErr := t.tx.Rollback(); rollbackErr != nil {
			*cause = fmt.Errorf("rollback caused by: %w, failed with: %v", *cause, rollbackErr)
		}
	}
}

func (t *TxController) Commit() error {
	if !t.isLocked {
		t.isLocked = true
		err := t.tx.Commit()
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	return nil
}

func (t *TxController) IsLocked() bool {
	return t.isLocked
}
