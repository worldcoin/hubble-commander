package api

import (
	"fmt"
	"testing"

	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
)

func TestSanitizeError_ValidateErrorsFromMap(t *testing.T) {
	sampleError1 := fmt.Errorf("sample error 1")
	sampleError2 := fmt.Errorf("sample error 2")
	sampleError3 := &storage.NotFoundError{}
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
		sampleError1: expectedAPIError1,
		sampleError2: expectedAPIError2,
		sampleError3: expectedAPIError3,
	}

	apiError := sanitizeError(sampleError1, errMap)
	require.Equal(t, *expectedAPIError1, *apiError)

	apiError = sanitizeError(sampleError2, errMap)
	require.Equal(t, *expectedAPIError2, *apiError)

	apiError = sanitizeError(sampleError3, errMap)
	require.Equal(t, *expectedAPIError3, *apiError)
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
	newErr3 := &storage.NotFoundError{}

	testCommonErrors := []*InternalToAPIError{
		NewInternalToAPIError(
			expectedAPIError1.Code,
			expectedAPIError1.Message,
			[]interface{}{
				newErr1,
				newErr2,
			},
		),
		NewInternalToAPIError(
			expectedAPIError2.Code,
			expectedAPIError2.Message,
			[]interface{}{
				newErr3,
			},
		),
	}

	apiError := sanitizeCommonError(newErr1, testCommonErrors)
	require.Equal(t, expectedAPIError1, *apiError)

	apiError = sanitizeCommonError(newErr2, testCommonErrors)
	require.Equal(t, expectedAPIError1, *apiError)

	apiError = sanitizeCommonError(newErr3, testCommonErrors)
	require.Equal(t, expectedAPIError2, *apiError)

	apiError = sanitizeCommonError(fmt.Errorf("ducks"), testCommonErrors)
	require.Equal(t, expectedUnknownError, *apiError)
}
