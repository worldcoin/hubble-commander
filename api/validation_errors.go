package api

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/pkg/errors"
)

var (
	AnyMissingFieldError, anyMissingFieldErrorSupport = utils.NewAnyError(&MissingFieldError{})
	AnyInvalidSignatureError                          = &InvalidSignatureError{}
)

type MissingFieldError struct {
	*utils.AnyErrorSupport
	field string
}

func NewMissingFieldError(field string) *MissingFieldError {
	return &MissingFieldError{
		AnyErrorSupport: anyMissingFieldErrorSupport,
		field:           field,
	}
}

func (e *MissingFieldError) Unwrap() error {
	return e.AnyErrorSupport
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("missing required %s field", e.field)
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
