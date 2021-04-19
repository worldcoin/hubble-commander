package api

import "fmt"

type MissingFieldError struct {
	field string
}

func NewMissingFieldError(field string) *MissingFieldError {
	return &MissingFieldError{field}
}

func (m MissingFieldError) Error() string {
	return fmt.Sprintf("missing required %s field", m.field)
}

type NotDecimalEncodableError struct {
	field string
}

func NewNotDecimalEncodableError(field string) *NotDecimalEncodableError {
	return &NotDecimalEncodableError{field: field}
}

func (e NotDecimalEncodableError) Error() string {
	return fmt.Sprintf("%s is not encodable as multi-precission decimal", e.field)
}
