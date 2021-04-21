package txstatus

type TransactionStatus uint

const (
	Pending TransactionStatus = iota + 1000
	InBatch
	Finalised                   // nolint:misspell
	Error     TransactionStatus = 5000
)

var TransactionStatuses = map[TransactionStatus]string{
	Pending:   "PENDING",
	InBatch:   "IN_BATCH",
	Finalised: "FINALISED", // nolint:misspell
	Error:     "ERROR",
}

func (s TransactionStatus) Ref() *TransactionStatus {
	return &s
}

func (s TransactionStatus) String() string {
	msg, exists := TransactionStatuses[s]
	if !exists {
		return "UNKNOWN"
	}
	return msg
}
