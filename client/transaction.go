package client

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/pkg/errors"
)

var (
	ErrMissingType    = errors.New("missing type")
	ErrNotImplemented = errors.New("not implemented")
)

type Transaction struct {
	Parsed models.GenericTransaction
}

func (tx *Transaction) UnmarshalJSON(bytes []byte) error {
	var rawTx struct {
		TxType *txtype.TransactionType
	}
	err := json.Unmarshal(bytes, &rawTx)
	if err != nil {
		return err
	}

	if rawTx.TxType == nil {
		return ErrMissingType
	}

	switch *rawTx.TxType {
	case txtype.Transfer:
		return tx.unmarshalTransfer(bytes)
	case txtype.Create2Transfer:
		return tx.unmarshalCreate2Transfer(bytes)
	case txtype.MassMigration:
		return tx.unmarshalMassMigration(bytes)
	default:
		return ErrNotImplemented
	}
}

func (tx *Transaction) unmarshalTransfer(bytes []byte) error {
	var transfer models.Transfer
	err := json.Unmarshal(bytes, &transfer)
	if err != nil {
		return err
	}
	tx.Parsed = &transfer
	return nil
}

func (tx *Transaction) unmarshalCreate2Transfer(bytes []byte) error {
	var transfer models.Create2Transfer
	err := json.Unmarshal(bytes, &transfer)
	if err != nil {
		return err
	}
	tx.Parsed = &transfer
	return nil
}

func (tx *Transaction) unmarshalMassMigration(bytes []byte) error {
	var transfer models.MassMigration
	err := json.Unmarshal(bytes, &transfer)
	if err != nil {
		return err
	}
	tx.Parsed = &transfer
	return nil
}

func txsToTransactionArray(txs []Transaction) models.GenericTransactionArray {
	genericTxs := make([]models.GenericTransaction, 0, len(txs))
	for i := range txs {
		genericTxs = append(genericTxs, txs[i].Parsed)
	}
	return models.MakeGenericArray(genericTxs...)
}
