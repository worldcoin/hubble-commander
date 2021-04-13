package dto

import (
	"encoding/json"
	"errors"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

var (
	ErrMissingType     = errors.New("missing type")
	ErrUnsupportedType = errors.New("unsupported type")
	ErrNotImplemented  = errors.New("not implemented")
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
		return ErrMissingType
	}

	switch *rawTx.Type {
	case txtype.Transfer:
		return tx.unmarshalTransfer(bytes)
	case txtype.Create2Transfer:
		return tx.unmarshalCreate2Transfer(bytes)
	case txtype.MassMigration:
		return ErrNotImplemented
	default:
		return ErrUnsupportedType
	}
}

func (tx *Transaction) unmarshalTransfer(bytes []byte) error {
	var transfer Transfer
	err := json.Unmarshal(bytes, &transfer)
	if err != nil {
		return err
	}
	tx.Parsed = transfer
	return nil
}

func (tx *Transaction) unmarshalCreate2Transfer(bytes []byte) error {
	var transfer Create2Transfer
	err := json.Unmarshal(bytes, &transfer)
	if err != nil {
		return err
	}
	tx.Parsed = transfer
	return nil
}
