package storage

import (
	"fmt"
)

type NotFoundErr struct {
	field string
}

func NewNotFoundError(field string) *NotFoundErr {
	return &NotFoundErr{field: field}
}

func (n *NotFoundErr) Error() string {
	return fmt.Sprintf("%s not found", n.field)
}
