package models

import "fmt"

type TransactionStatus uint

const (
	Pending   TransactionStatus = 1001
	Committed TransactionStatus = 1002
	InBatch   TransactionStatus = 1003
	Finalized TransactionStatus = 1004
	Error     TransactionStatus = 5000
)

var TransactionsStatuses = [5]TransactionStatus{
	Pending,
	Committed,
	InBatch,
	Finalized,
	Error,
}

func (s TransactionStatus) String() string {
	return fmt.Sprintf("%d", s)
}

func (s TransactionStatus) Message() string {
	switch s {
	case Pending:
		return "PENDING"

	case Committed:
		return "COMMITTED"

	case InBatch:
		return "IN_BATCH"

	case Finalized:
		return "FINALIZED"

	case Error:
		return "ERROR"

	default:
		return "UNKNOWN"
	}
}
