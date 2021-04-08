package dto

import (
	"encoding/json"
	"fmt"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

type Transaction struct {
	Parsed interface{}
}

func MakeTransaction(parsed interface{}) Transaction {
	return Transaction{Parsed: parsed}
}

func (tx *Transaction) UnmarshalJSON(bytes []byte) error {
	var rawTx struct {
		Type *txtype.TransactionType
	}
	err := json.Unmarshal(bytes, &rawTx)
	if err != nil {
		return err
	}

	if rawTx.Type == nil {
		return fmt.Errorf("missing type")
	}

	switch *rawTx.Type {
	case txtype.Transfer:
		tx.Parsed = Transfer{}
		return json.Unmarshal(bytes, tx.Parsed)
	default:
		return fmt.Errorf("unsupported type")
	}
}
