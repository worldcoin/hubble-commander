package api

import (
	"errors"
	"fmt"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/db"
)

var usedErrorCodes = map[int]bool{}

const unknownAPIError = 999

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

func (e *NotDecimalEncodableError) Is(other error) bool {
	otherError, ok := other.(*NotDecimalEncodableError)
	if !ok {
		return false
	}
	return *e == *otherError
}

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("error code: %d", e.Code)
	}
	return e.Message
}

func (e *APIError) ErrorCode() int {
	return e.Code
}

func NewAPIError(code int, message string) *APIError {
	if usedErrorCodes[code] {
		panic(fmt.Sprintf("%d API error code is already used", code))
	}

	usedErrorCodes[code] = true

	return &APIError{
		Code:    code,
		Message: message,
	}
}

func NewUnknownAPIError(err error) *APIError {
	unknownAPIErrorMessage := fmt.Sprintf("unknown error: %s", err.Error())

	if usedErrorCodes[unknownAPIError] {
		return &APIError{
			Code:    unknownAPIError,
			Message: unknownAPIErrorMessage,
		}
	}

	return NewAPIError(unknownAPIError, fmt.Sprintf("unknown error: %s", err.Error()))
}

type InternalToAPIError struct {
	apiError        *APIError
	commanderErrors []error
}

func NewInternalToAPIError(code int, message string, commanderErrors []error) *InternalToAPIError {
	return &InternalToAPIError{
		apiError:        NewAPIError(code, message),
		commanderErrors: commanderErrors,
	}
}

var commonErrors = []*InternalToAPIError{
	// Badger
	NewInternalToAPIError(
		40000,
		"an error occurred while saving data to the Badger database",
		[]error{
			db.ErrInconsistentItemsLength,
			db.ErrInvalidKeyListLength,
			db.ErrInvalidKeyListMetadataLength,
		},
	),
	NewInternalToAPIError(40001, "an error occurred while iterating over Badger database", []error{db.ErrIteratorFinished}),
	// BLS
	NewInternalToAPIError(99004, "an error occurred while fetching the domain for signing", []error{bls.ErrInvalidDomainLength}),
}

func sanitizeError(err error, errMap map[error]*APIError) *APIError {
	for k, v := range errMap {
		if errors.Is(err, k) {
			return v
		}
	}

	return sanitizeCommonError(err, commonErrors)
}

func sanitizeCommonError(err error, errMap []*InternalToAPIError) *APIError {
	for i := range errMap {
		selectedErrMap := errMap[i]
		for j := range selectedErrMap.commanderErrors {
			commanderErr := selectedErrMap.commanderErrors[j]
			if errors.Is(err, commanderErr) {
				return errMap[i].apiError
			}
		}
	}

	return NewUnknownAPIError(err)
}
