package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Transfer struct {
	TransactionBase
	ToStateID uint32 `db:"to_state_id"`
}

type TransferForCommitment struct {
	TransactionBaseForCommitment
	ToStateID uint32 `db:"to_state_id"`
}

type TransferWithBatchHash struct {
	Transfer
	BatchHash *common.Hash `db:"batch_hash"`
}

func (t *Transfer) Type() txtype.TransactionType {
	return txtype.Transfer
}

func (t *Transfer) GetBase() *TransactionBase {
	return &t.TransactionBase
}

func (t *Transfer) GetToStateID() *uint32 {
	return &t.ToStateID
}

// nolint:gocritic
func (t Transfer) Copy() GenericTransaction {
	return &t
}
