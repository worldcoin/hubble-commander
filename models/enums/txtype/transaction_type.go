package txtype

import (
	"encoding/json"

	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type TransactionType uint8

const (
	Genesis         TransactionType = 0
	Transfer        TransactionType = 1
	MassMigration   TransactionType = 2
	Create2Transfer TransactionType = 3
	Deposit         TransactionType = 4
)

var TransactionTypes = map[TransactionType]string{
	Genesis:         "GENESIS",
	Transfer:        "TRANSFER",
	Create2Transfer: "CREATE2TRANSFER",
	MassMigration:   "MASS_MIGRATION",
	Deposit:         "DEPOSIT",
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
