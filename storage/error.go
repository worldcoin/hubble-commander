package storage

import (
	"errors"
	"fmt"
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

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", n.field)
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	target := &NotFoundError{}
	return errors.As(err, &target)
}
