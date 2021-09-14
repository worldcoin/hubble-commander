package api

import (
	"errors"
	"fmt"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/db"
)

var usedErrorCodes map[int]bool

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
	if usedErrorCodes == nil {
		usedErrorCodes = map[int]bool{}
	}

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
	return NewAPIError(999, fmt.Sprintf("unknown error: %s", err.Error()))
}

type CommanderErrorsToErrorAPI struct {
	apiError        *APIError
	commanderErrors []interface{}
}

func NewCommanderErrorsToErrorAPI(code int, message string, commanderErrors []interface{}) *CommanderErrorsToErrorAPI {
	return &CommanderErrorsToErrorAPI{
		apiError:        NewAPIError(code, message),
		commanderErrors: commanderErrors,
	}
}

var commonErrors = []*CommanderErrorsToErrorAPI{
	// Badger
	NewCommanderErrorsToErrorAPI(
		40000,
		"an error occurred while saving data to the Badger database",
		[]interface{}{
			db.ErrInconsistentItemsLength,
			db.ErrInvalidKeyListLength,
			db.ErrInvalidKeyListMetadataLength,
		},
	),
	NewCommanderErrorsToErrorAPI(40001, "an error occurred while iterating over badger database", []interface{}{db.ErrIteratorFinished}),
	// BLS
	NewCommanderErrorsToErrorAPI(99004, "an error occurred while fetching the domain for signing", []interface{}{bls.ErrInvalidDomainLength}),
}

func sanitizeError(err error, errMap map[error]*APIError) *APIError {
	for k, v := range errMap {
		if errors.Is(err, k) {
			return v
		}
	}

	return sanitizeCommonError(err, commonErrors)
}

func sanitizeCommonError(err error, errMap []*CommanderErrorsToErrorAPI) *APIError {
	for i := range errMap {
		selectedErrMap := errMap[i]
		for j := range selectedErrMap.commanderErrors {
			commanderErr := selectedErrMap.commanderErrors[j]
			if errors.Is(err, commanderErr.(error)) {
				return errMap[i].apiError
			}
		}
	}

	return NewUnknownAPIError(err)
}
