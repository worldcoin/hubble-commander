package api

import "fmt"

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

type ErrorAPI struct {
	Code    int
	Message string
}

func (e *ErrorAPI) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("error code: %d", e.Code)
	}
	return e.Message
}

func (e *ErrorAPI) ErrorCode() int {
	return e.Code
}

// TODO - verify that there are no duplicates here
func NewAPIError(code int, message string) *ErrorAPI {
	return &ErrorAPI{
		Code:    code,
		Message: message,
	}
}

type CommanderErrorsToErrorAPI struct {
	apiError *ErrorAPI
	//commanderErrors []error
	commanderErrors []interface{} // temp
}

func NewCommanderErrorsToErrorAPI(code int, message string, commanderErrors []interface{}) *CommanderErrorsToErrorAPI {
	return &CommanderErrorsToErrorAPI{
		apiError:        NewAPIError(code, message),
		commanderErrors: commanderErrors,
	}
}

var unknownError = ErrorAPI{
	Code:    999,
	Message: "unknown error",
}

/*
	ERROR CODES:
	10XXX - Transactions
	20XXX - Commitments
	30XXX - Batches
	40XXX - Badger?


*/

var commonErrors = []*CommanderErrorsToErrorAPI{
	// Badger
	NewCommanderErrorsToErrorAPI(
		40000,
		"an error occurred while saving data to the Badger database",
		[]interface{}{
			badger.ErrInconsistentItemsLength,
			badger.ErrInvalidKeyListLength,
			badger.ErrInvalidKeyListMetadataLength,
		},
	),

	// Send transactions
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{encoder.ErrNotEncodableDecimal}),

	// handle transfer
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{MissingFieldError{}}),
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{ErrInvalidAmount}),
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{ErrFeeTooLow}),
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{ErrTransferToSelf}),
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{storage.NotFoundError{}}),
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{ErrNonceTooLow}),
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{ErrNonceTooHigh}),
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{ErrNotEnoughBalance}),
	NewCommanderErrorsToErrorAPI(0, "", []interface{}{storage.NotFoundError{"account leaf"}}), // TODO-API decide how to handle these errors
	NewCommanderErrorsToErrorAPI(0, "signing error", []interface{}{bls.ErrInvalidDomainLength, ErrInvalidSignature}),
	// how do I handle bls signing errors?
	// add pack errors
}

func sanitizeError(err error, x map[error]ErrorAPI) *ErrorAPI {
	for k, v := range x {
		if errors.As(err, &k) {
			return &v
		}
	}

	return sanitizeCommonError(err, commonErrors)
}

func sanitizeCommonError(err error, x map[error]ErrorAPI) *ErrorAPI {
	for k, v := range x {
		if errors.As(err, &k) {
			return &v
		}
	}

	return &unknownError
}
