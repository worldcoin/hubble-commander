package txstatus

import (
	"encoding/json"

	cs "github.com/Worldcoin/hubble-commander/models/enums/commitmentstatus"
	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type TransactionStatus uint

const (
	Pending                     = TransactionStatus(cs.Pending)
	Mined                       = TransactionStatus(cs.Mined)
	Finalised                   = TransactionStatus(cs.Finalised) // nolint:misspell
	Error     TransactionStatus = 5000
)

var TransactionStatuses = map[TransactionStatus]string{
	Pending:   cs.CommitmentStatuses[cs.Pending],
	Mined:     cs.CommitmentStatuses[cs.Mined],
	Finalised: cs.CommitmentStatuses[cs.Finalised], // nolint:misspell
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
