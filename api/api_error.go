package api

import (
	"fmt"
)

var usedErrorCodes = map[int]bool{}

const unknownAPIErrorCode = 999

type APIError struct {
	Code    int
	Message string
	Data    interface{} `json:",omitempty"`
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

func (e *APIError) ErrorData() interface{} {
	return e.Data
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

	if usedErrorCodes[unknownAPIErrorCode] {
		return &APIError{
			Code:    unknownAPIErrorCode,
			Message: unknownAPIErrorMessage,
		}
	}

	return NewAPIError(unknownAPIErrorCode, fmt.Sprintf("unknown error: %s", err.Error()))
}
