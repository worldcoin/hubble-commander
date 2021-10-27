package api

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	AnyMissingFieldError     = &MissingFieldError{field: anythingField}
	AnyInvalidSignatureError = &InvalidSignatureError{}
)

const anythingField = "anything"

type MissingFieldError struct {
	field string
}

func NewMissingFieldError(field string) *MissingFieldError {
	if field == anythingField {
		panic(fmt.Sprintf(`cannot use "%s" field for MissingFieldError`, anythingField))
	}
	return &MissingFieldError{field}
}

func (m *MissingFieldError) Error() string {
	return fmt.Sprintf("missing required %s field", m.field)
}

func (m *MissingFieldError) Is(other error) bool {
	otherError, ok := other.(*MissingFieldError)
	if !ok {
		return false
	}
	if *m == *AnyMissingFieldError || *otherError == *AnyMissingFieldError {
		return true
	}
	return *m == *otherError
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

func (e *NotDecimalEncodableError) Is(other error) bool {
	otherError, ok := other.(*NotDecimalEncodableError)
	if !ok {
		return false
	}
	return *e == *otherError
}

type InvalidSignatureError struct {
	reason string
}

func NewInvalidSignatureError(reason string) *InvalidSignatureError {
	return &InvalidSignatureError{reason: reason}
}

func (e InvalidSignatureError) Error() string {
	return fmt.Sprintf("invalid signature: %s", e.reason)
}

func (e *InvalidSignatureError) Is(other error) bool {
	var invalidSignatureErr *InvalidSignatureError
	return errors.As(other, &invalidSignatureErr)
}
