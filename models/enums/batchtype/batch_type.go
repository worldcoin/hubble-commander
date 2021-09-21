package batchtype

import (
	"encoding/json"

	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type BatchType uint8

const (
	Genesis BatchType = iota
	Transfer
	MassMigration
	Create2Transfer
	Deposit
)

var BatchTypes = map[BatchType]string{
	Genesis:         "GENESIS",
	Transfer:        "TRANSFER",
	MassMigration:   "MASS_MIGRATION",
	Create2Transfer: "CREATE2TRANSFER",
	Deposit:         "DEPOSIT",
}

func (s BatchType) Ref() *BatchType {
	return &s
}

func (s BatchType) String() string {
	msg, exists := BatchTypes[s]
	if !exists {
		return "UNKNOWN"
	}
	return msg
}

func (s *BatchType) UnmarshalJSON(bytes []byte) error {
	var strType string
	err := json.Unmarshal(bytes, &strType)
	if err != nil {
		return err
	}

	for k, v := range BatchTypes {
		if v == strType {
			*s = k
			return nil
		}
	}
	return enumerr.NewUnsupportedError("batch type")
}

func (s BatchType) MarshalJSON() ([]byte, error) {
	msg, exists := BatchTypes[s]
	if !exists {
		return nil, enumerr.NewUnsupportedError("batch type")
	}
	return json.Marshal(msg)
}
