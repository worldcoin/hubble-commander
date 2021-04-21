package txtype

import (
	"encoding/json"
	"errors"
)

type TransactionType uint8

const (
	Transfer        TransactionType = 1
	Create2Transfer TransactionType = 3
	MassMigration   TransactionType = 5
)

var (
	TransactionTypes = map[TransactionType]string{
		Transfer:        "TRANSFER",
		Create2Transfer: "CREATE2TRANSFER",
		MassMigration:   "MASS_MIGRATION",
	}

	ErrUnsupportedTransactionType = errors.New("unsupported transaction type")
)

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
	return ErrUnsupportedTransactionType
}

func (s TransactionType) MarshalJSON() ([]byte, error) {
	msg, exists := TransactionTypes[s]
	if !exists {
		return nil, ErrUnsupportedTransactionType
	}
	return json.Marshal(msg)
}
