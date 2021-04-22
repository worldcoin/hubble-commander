package errors

import (
	"errors"
	"fmt"
)

type UnsupportedError struct {
	field string
}

func NewUnsupportedError(field string) *UnsupportedError {
	return &UnsupportedError{field: field}
}

func (n *UnsupportedError) Error() string {
	return fmt.Sprintf("unsupported %s", n.field)
}

func IsUnsupportedError(err error) bool {
	if err == nil {
		return false
	}
	target := &UnsupportedError{}
	return errors.As(err, &target)
}
