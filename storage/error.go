package storage

import (
	"errors"
	"fmt"

	err "github.com/pkg/errors"
)

type NotFoundError struct {
	field string
}

func NewNotFoundError(field string) error {
	return err.WithStack(&NotFoundError{field: field})
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
