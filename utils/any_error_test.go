package utils

import (
	stdErr "errors"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var anyNotFoundError, anyNotFoundErrorSupport = NewAnyError(&notFoundError{})

type notFoundError struct {
	*AnyErrorSupport
	field string
}

func newNotFoundError(field string) *notFoundError {
	return &notFoundError{
		AnyErrorSupport: anyNotFoundErrorSupport,
		field:           field,
	}
}

func (e *notFoundError) Unwrap() error {
	return e.AnyErrorSupport
}

func (e *notFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.field)
}

func TestAnyErrorSupport_DoesNotAffectDefaultIsComparison(t *testing.T) {
	notFound := newNotFoundError("something")
	notFoundCopy := notFound
	require.ErrorIs(t, notFound, notFoundCopy)
	require.ErrorIs(t, errors.WithStack(notFound), notFoundCopy)

	randomError := stdErr.New("random error")
	require.NotErrorIs(t, notFound, randomError)

	otherNotFound := newNotFoundError("something else")
	require.NotErrorIs(t, notFound, otherNotFound)
}

func TestAnyError_WorksWithDefaultIsComparison(t *testing.T) {
	notFound := newNotFoundError("something")
	require.ErrorIs(t, notFound, anyNotFoundError)
	require.ErrorIs(t, errors.WithStack(notFound), anyNotFoundError)

	randomError := stdErr.New("random error")
	require.NotErrorIs(t, randomError, anyNotFoundError)
}

type notFoundErrorWithIs struct {
	notFoundError
	ignoredField string
}

func newNotFoundErrorWithIs(field, ignoredField string) *notFoundErrorWithIs {
	return &notFoundErrorWithIs{
		notFoundError: *newNotFoundError(field),
		ignoredField:  ignoredField,
	}
}

func (a *notFoundErrorWithIs) Is(other error) bool {
	otherError, ok := other.(*notFoundErrorWithIs)
	if !ok {
		return false
	}
	return a.field == otherError.field
}

func TestAnyErrorSupport_DoesNotAffectOverriddenIsComparison(t *testing.T) {
	notFound := newNotFoundErrorWithIs("some value", "some ignored value")
	notFoundCopy := notFound
	notFoundCopy.ignoredField = "other ignored value"
	require.ErrorIs(t, notFound, notFoundCopy)

	randomError := stdErr.New("random error")
	require.NotErrorIs(t, notFound, randomError)

	otherNotFound := newNotFoundErrorWithIs("other value", "some ignored value")
	require.NotErrorIs(t, notFound, otherNotFound)
}

func TestAnyError_WorksWithOverriddenIsComparison(t *testing.T) {
	notFound := newNotFoundErrorWithIs("some value", "some ignored value")
	require.ErrorIs(t, notFound, anyNotFoundError)
	require.ErrorIs(t, errors.WithStack(notFound), anyNotFoundError)
}

var anyInvalidArgumentError, anyInvalidArgumentErrorSupport = NewAnyError(&invalidArgumentError{})

type invalidArgumentError struct {
	*AnyErrorSupport
	arg string
}

func newInvalidArgumentError(arg string) *invalidArgumentError {
	return &invalidArgumentError{
		AnyErrorSupport: anyInvalidArgumentErrorSupport,
		arg:             arg,
	}
}

func (e *invalidArgumentError) Unwrap() error {
	return e.AnyErrorSupport
}

func TestAnyError_TakesIntoAccountErrorType(t *testing.T) {
	invalidArg := newInvalidArgumentError("something")
	require.ErrorIs(t, invalidArg, anyInvalidArgumentError)
	require.NotErrorIs(t, invalidArg, anyNotFoundError)
}
