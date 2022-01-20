package txstatus

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type TransactionStatus uint

const (
	Pending                     = TransactionStatus(batchstatus.Pending)
	InBatch                     = TransactionStatus(batchstatus.InBatch)
	Finalised                   = TransactionStatus(batchstatus.Finalised)
	Error     TransactionStatus = 5000
)

var TransactionStatuses = map[TransactionStatus]string{
	Pending:   batchstatus.BatchStatuses[batchstatus.Pending],
	InBatch:   batchstatus.BatchStatuses[batchstatus.InBatch],
	Finalised: batchstatus.BatchStatuses[batchstatus.Finalised],
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
	return enumerr.NewUnsupportedError("transaction status")
}

func (s TransactionStatus) MarshalJSON() ([]byte, error) {
	msg, exists := TransactionStatuses[s]
	if !exists {
		return nil, enumerr.NewUnsupportedError("transaction status")
	}
	return json.Marshal(msg)
}
