package storage

import (
	"errors"
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
)

var (
	ErrNoRowsAffected          = fmt.Errorf("no rows were affected by the update")
	ErrNonexistentState        = fmt.Errorf("cannot revert to nonexistent state")
	ErrAlreadyMinedTransaction = fmt.Errorf("transaction already mined")

	AnyNotFoundError, anyNotFoundErrorSupport = utils.NewAnyError(&NotFoundError{})
)

type NotFoundError struct {
	*utils.AnyErrorSupport
	field string
}

func NewNotFoundError(field string) *NotFoundError {
	return &NotFoundError{
		AnyErrorSupport: anyNotFoundErrorSupport,
		field:           field,
	}
}

func (e *NotFoundError) Unwrap() error {
	return e.AnyErrorSupport
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.field)
}

func (e *NotFoundError) Is(other error) bool {
	otherError, ok := other.(*NotFoundError)
	if !ok {
		return false
	}
	return *e == *otherError
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

type NoVacantSubtreeError struct {
	subtreeDepth uint8
}

func NewNoVacantSubtreeError(subtreeDepth uint8) *NoVacantSubtreeError {
	return &NoVacantSubtreeError{subtreeDepth: subtreeDepth}
}

func (e *NoVacantSubtreeError) Error() string {
	return fmt.Sprintf("no vacant slot found in the State Tree for a subtree of depth %d", e.subtreeDepth)
}

func (e *NoVacantSubtreeError) Is(other error) bool {
	otherError, ok := other.(*NoVacantSubtreeError)
	if !ok {
		return false
	}
	return *e == *otherError
}
