package models

import "github.com/ethereum/go-ethereum/common"

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

func (t *Transfer) GetToStateID() *uint32 {
	return &t.ToStateID
}

func (t Transfer) Copy() GenericTransfer {
	return &t
}
