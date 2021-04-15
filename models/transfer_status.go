package models

type TransferStatus uint

const (
	Pending TransferStatus = iota + 1000
	Committed
	InBatch
	Finalised                   // nolint:misspell
	Error     TransferStatus = 5000
)

var TransferStatuses = map[TransferStatus]string{
	Pending:   "PENDING",
	Committed: "COMMITTED",
	InBatch:   "IN_BATCH",
	Finalised: "FINALISED", // nolint:misspell
	Error:     "ERROR",
}

func (s TransferStatus) Ref() *TransferStatus {
	return &s
}

func (s TransferStatus) String() string {
	msg, exists := TransferStatuses[s]
	if !exists {
		return "UNKNOWN"
	}
	return msg
}
