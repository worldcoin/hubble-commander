package storage

import (
	"errors"
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrNoRowsAffected   = errors.New("no rows were affected by the update")
	ErrNotExistentState = errors.New("cannot revert to not existent state")
)

type NotFoundError struct {
	field string
}

func NewNotFoundError(field string) *NotFoundError {
	return &NotFoundError{field: field}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.field)
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	target := &NotFoundError{}
	return errors.As(err, &target)
}

type InvalidPubKeyIDError struct {
	value uint32
}

func NewInvalidPubKeyIDError(value uint32) *InvalidPubKeyIDError {
	return &InvalidPubKeyIDError{value: value}
}

func (e *InvalidPubKeyIDError) Error() string {
	return fmt.Sprintf("invalid pubKeyID value: %d", e.value)
}

type AccountAlreadyExistsError struct {
	Account *models.AccountLeaf
}

func NewAccountAlreadyExistsError(account *models.AccountLeaf) *AccountAlreadyExistsError {
	return &AccountAlreadyExistsError{Account: account}
}

func (e *AccountAlreadyExistsError) Error() string {
	return fmt.Sprintf("account with pubKeyID %d already exists", e.Account.PubKeyID)
}

type AccountBatchAlreadyExistsError struct {
	Accounts []models.AccountLeaf
}

func NewAccountBatchAlreadyExistsError(accounts []models.AccountLeaf) *AccountBatchAlreadyExistsError {
	return &AccountBatchAlreadyExistsError{Accounts: accounts}
}

func (e *AccountBatchAlreadyExistsError) Error() string {
	return fmt.Sprintf("accounts with pubKeyIDs %v already exist", e.Accounts)
}
