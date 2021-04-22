package models

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionBase struct {
	Hash                 common.Hash            `db:"tx_hash"`
	TxType               txtype.TransactionType `db:"tx_type"`
	FromStateID          uint32                 `db:"from_state_id"`
	Amount               Uint256
	Fee                  Uint256
	Nonce                Uint256
	Signature            []byte
	IncludedInCommitment *int32  `db:"included_in_commitment"`
	ErrorMessage         *string `db:"error_message"`
}

type Transaction struct {
	Parsed interface{}
}

func (t *Transaction) Scan(src interface{}) error {
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

	switch *rawTx.Type { // nolint:exhaustive
	case txtype.Transfer:
		return tx.unmarshalTransfer(bytes)
	case txtype.Create2Transfer:
		return tx.unmarshalCreate2Transfer(bytes)
	default:
		return ErrNotImplemented
	}
}
