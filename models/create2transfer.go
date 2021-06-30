package models

import "github.com/ethereum/go-ethereum/common"

type Create2Transfer struct {
	TransactionBase
	ToStateID   *uint32   `db:"to_state_id"`
	ToPublicKey PublicKey `db:"to_public_key"`
}

type Create2TransferForCommitment struct {
	TransactionBaseForCommitment
	ToStateID   *uint32   `db:"to_state_id"`
	ToPublicKey PublicKey `db:"to_public_key"`
}

type Create2TransferWithBatchHash struct {
	Create2Transfer
	BatchHash *common.Hash `db:"batch_hash"`
}

func (t *Create2Transfer) GetToStateID() *uint32 {
	return t.ToStateID
}

func (t Create2Transfer) Copy() GenericTransfer {
	return &t
}
