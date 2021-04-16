package storage

import (
	"errors"
	"fmt"
)

type NotFoundErr struct {
	field string
}

func NewNotFoundErr(field string) *NotFoundErr {
	return &NotFoundErr{field: field}
}

func (n *NotFoundErr) Error() string {
	return fmt.Sprintf("%s not found", n.field)
}

func IsNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	target := &NotFoundErr{}
	return errors.As(err, &target)
}
