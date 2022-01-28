package txstatus

import (
	"encoding/json"

	bs "github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type TransactionStatus uint

const (
	Pending   TransactionStatus = 999
	Submitted                   = TransactionStatus(bs.Submitted)
	Mined                       = TransactionStatus(bs.Mined)
	Finalised                   = TransactionStatus(bs.Finalised) // nolint:misspell
	Error     TransactionStatus = 5000
)

var TransactionStatuses = map[TransactionStatus]string{
	Pending:   "PENDING",
	Submitted: bs.BatchStatuses[bs.Submitted],
	Mined:     bs.BatchStatuses[bs.Mined],
	Finalised: bs.BatchStatuses[bs.Finalised], // nolint:misspell
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
