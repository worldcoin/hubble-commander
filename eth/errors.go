package eth

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/utils/ref"
)

type LogNotFoundError struct {
	logName string
}

func NewLogNotFoundError(logName string) *LogNotFoundError {
	return &LogNotFoundError{logName}
}

func (e LogNotFoundError) Error() string {
	return fmt.Sprintf("log not found in the receipt: %s", e.logName)
}

type DisputeTxRevertedError struct {
	batchID uint64
	reason  *string
}

func NewDisputeTxRevertedError(batchID uint64, reason string) *DisputeTxRevertedError {
	return &DisputeTxRevertedError{
		batchID: batchID,
		reason:  ref.String(reason),
	}
}

func NewUnknownDisputeTxRevertedError(batchID uint64) *DisputeTxRevertedError {
	return &DisputeTxRevertedError{
		batchID: batchID,
	}
}

func (e DisputeTxRevertedError) Error() string {
	msg := fmt.Sprintf("dispute of batch #%d failed", e.batchID)
	if e.reason != nil {
		msg += fmt.Sprintf(": %s", *e.reason)
	}
	return msg
}
