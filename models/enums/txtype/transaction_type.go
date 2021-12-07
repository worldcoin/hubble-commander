package txtype

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type TransactionType uint8

const (
	Transfer        = TransactionType(batchtype.Transfer)
	Create2Transfer = TransactionType(batchtype.Create2Transfer)
	MassMigration   = 5
)

var TransactionTypes = map[TransactionType]string{
	Transfer:        "TRANSFER",
	Create2Transfer: "CREATE2TRANSFER",
	MassMigration:   "MASS_MIGRATION",
}

func (s TransactionType) Ref() *TransactionType {
	return &s
}

func (s TransactionType) String() string {
	msg, exists := TransactionTypes[s]
	if !exists {
		return "UNKNOWN"
	}
	return msg
}

func (s *TransactionType) UnmarshalJSON(bytes []byte) error {
	var strType string
	err := json.Unmarshal(bytes, &strType)
	if err != nil {
		return err
	}

	for k, v := range TransactionTypes {
		if v == strType {
			*s = k
			return nil
		}
	}
	return enumerr.NewUnsupportedError("transaction type")
}

func (s TransactionType) MarshalJSON() ([]byte, error) {
	msg, exists := TransactionTypes[s]
	if !exists {
		return nil, enumerr.NewUnsupportedError("transaction type")
	}
	return json.Marshal(msg)
}

func FromBatchType(batchType batchtype.BatchType) TransactionType {
	if batchType == batchtype.MassMigration {
		return MassMigration
	}
	return TransactionType(batchType)
}
