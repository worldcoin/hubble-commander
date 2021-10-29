package api

import (
	"fmt"
	"testing"

	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestSanitizeError_ValidateErrorsFromMap(t *testing.T) {
	sampleError1 := fmt.Errorf("sample error 1")
	sampleError2 := fmt.Errorf("sample error 2")
	expectedAPIError1 := &APIError{
		Code:    123,
		Message: "api error 1",
	}
	expectedAPIError2 := &APIError{
		Code:    234,
		Message: "api error 2",
	}
	expectedAPIError3 := &APIError{
		Code:    345,
		Message: "api error 3",
	}

	errMap := map[error]*APIError{
		sampleError1:             expectedAPIError1,
		sampleError2:             expectedAPIError2,
		storage.AnyNotFoundError: expectedAPIError3,
	}

	apiError := sanitizeError(sampleError1, errMap)
	require.Equal(t, *expectedAPIError1, *apiError)

	apiError = sanitizeError(errors.WithStack(sampleError1), errMap)
	require.Equal(t, *expectedAPIError1, *apiError)

	apiError = sanitizeError(sampleError2, errMap)
	require.Equal(t, *expectedAPIError2, *apiError)

	apiError = sanitizeError(errors.WithStack(sampleError2), errMap)
	require.Equal(t, *expectedAPIError2, *apiError)

	apiError = sanitizeError(storage.NewNotFoundError("something"), errMap)
	require.Equal(t, *expectedAPIError3, *apiError)

	apiError = sanitizeError(errors.WithStack(storage.NewNotFoundError("something")), errMap)
	require.Equal(t, *expectedAPIError3, *apiError)
}

func TestSanitizeError_ReturnsStackTraceInDebugLogLevel(t *testing.T) {
	expectedAPIError := &APIError{
		Code:    123,
		Message: "api error",
	}

	errMap := map[error]*APIError{
		storage.AnyNotFoundError: expectedAPIError,
	}

	apiError := sanitizeError(storage.NewNotFoundError("something"), errMap)
	require.Nil(t, apiError.Data)

	currentLevel := log.GetLevel()
	log.SetLevel(log.DebugLevel)
	defer func() {
		log.SetLevel(currentLevel)
	}()

	apiError = sanitizeError(storage.NewNotFoundError("something"), errMap)
	require.NotNil(t, apiError.Data)
}

func TestSanitizeCommonError(t *testing.T) {
	expectedAPIError1 := APIError{
		Code:    12345,
		Message: "expected error message 1",
	}
	expectedAPIError2 := APIError{
		Code:    54321,
		Message: "expected error message 2",
	}
	expectedUnknownError := APIError{
		Code:    999,
		Message: "unknown error: ducks",
	}
	newErr1 := fmt.Errorf("error1")
	newErr2 := fmt.Errorf("error2")
	testCommonErrors := []*InternalToAPIError{
		NewInternalToAPIError(
			expectedAPIError1.Code,
			expectedAPIError1.Message,
			[]error{
				newErr1,
				newErr2,
			},
		),
		NewInternalToAPIError(
			expectedAPIError2.Code,
			expectedAPIError2.Message,
			[]error{
				storage.AnyNotFoundError,
			},
		),
	}

	apiError := sanitizeCommonError(newErr1, testCommonErrors)
	require.Equal(t, expectedAPIError1, *apiError)

	apiError = sanitizeCommonError(errors.WithStack(newErr1), testCommonErrors)
	require.Equal(t, expectedAPIError1, *apiError)

	apiError = sanitizeCommonError(newErr2, testCommonErrors)
	require.Equal(t, expectedAPIError1, *apiError)

	apiError = sanitizeCommonError(errors.WithStack(newErr2), testCommonErrors)
	require.Equal(t, expectedAPIError1, *apiError)

	apiError = sanitizeCommonError(storage.NewNotFoundError("something"), testCommonErrors)
	require.Equal(t, expectedAPIError2, *apiError)

	apiError = sanitizeCommonError(errors.WithStack(storage.NewNotFoundError("something")), testCommonErrors)
	require.Equal(t, expectedAPIError2, *apiError)

	apiError = sanitizeCommonError(fmt.Errorf("ducks"), testCommonErrors)
	require.Equal(t, expectedUnknownError, *apiError)
}
