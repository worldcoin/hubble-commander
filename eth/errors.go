package eth

import "fmt"

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
}

func NewDisputeTxRevertedError(batchID uint64) *DisputeTxRevertedError {
	return &DisputeTxRevertedError{
		batchID: batchID,
	}
}

func (e DisputeTxRevertedError) Error() string {
	return fmt.Sprintf("dispute of batch #%d failed", e.batchID)
}
