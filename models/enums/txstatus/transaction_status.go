package txstatus

import (
	"encoding/json"
	"errors"
)

type TransactionStatus uint

const (
	Pending TransactionStatus = iota + 1000
	InBatch
	Finalised                   // nolint:misspell
	Error     TransactionStatus = 5000
)

var (
	TransactionStatuses = map[TransactionStatus]string{
		Pending:   "PENDING",
		InBatch:   "IN_BATCH",
		Finalised: "FINALISED", // nolint:misspell
		Error:     "ERROR",
	}

	ErrUnsupportedTransactionStatus = errors.New("unsupported transaction status")
)

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

func (s *TransactionStatus) UnmarshalJSON(bytes []byte) error {
	var strType string
	err := json.Unmarshal(bytes, &strType)
	if err != nil {
		return err
	}

	for k, v := range TransactionStatuses {
		if v == strType {
			*s = k
			return nil
		}
	}
	return ErrUnsupportedTransactionStatus
}

func (s TransactionStatus) MarshalJSON() ([]byte, error) {
	msg, exists := TransactionStatuses[s]
	if !exists {
		return nil, ErrUnsupportedTransactionStatus
	}
	return json.Marshal(msg)
}
