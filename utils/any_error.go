package utils

import (
	"fmt"
	"reflect"
)

// ErrorThatSupportsAny is an error that wraps AnyErrorSupport instance, created with NewAnyError.
// See the tests for this file for example usages.
type ErrorThatSupportsAny interface {
	Unwrap() error // should return pointer to the wrapped AnyErrorSupport struct

	// Left just as an additional check, AnyErrorSupport implements these interfaces
	error
	Is(other error) bool
}

type AnyError struct {
	errorType string
}

func NewAnyError(err ErrorThatSupportsAny) (*AnyError, *AnyErrorSupport) {
	errorType := reflect.ValueOf(err).Elem().Type().Name()
	return &AnyError{errorType}, &AnyErrorSupport{errorType}
}

func (e *AnyError) Error() string {
	return fmt.Sprintf("any %s", e.errorType)
}

type AnyErrorSupport struct {
	errorType string
}

func (a *AnyErrorSupport) Error() string {
	return fmt.Sprintf("error of %s type", a.errorType)
}

func (a *AnyErrorSupport) Is(other error) bool {
	otherError, ok := other.(*AnyError)
	if !ok {
		return false
	}
	return a.errorType == otherError.errorType
}
