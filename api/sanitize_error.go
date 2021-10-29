package api

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

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
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.Debugf("Sanitizing error:\n%+v", err)
	}

	for k, v := range errMap {
		if errors.Is(err, k) {
			if logrus.IsLevelEnabled(logrus.DebugLevel) {
				v.Data = fmt.Sprintf("%+v", err)
			}
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
