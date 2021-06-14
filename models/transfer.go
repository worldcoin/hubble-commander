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

func (t *Transfer) GetFromStateID() uint32 {
	return t.FromStateID
}

func (t *Transfer) GetToStateID() *uint32 {
	return &t.ToStateID
}

func (t *Transfer) GetAmount() Uint256 {
	return t.Amount
}

func (t *Transfer) GetFee() Uint256 {
	return t.Fee
}

func (t *Transfer) GetNonce() Uint256 {
	return t.Nonce
}

func (t *Transfer) SetNonce(nonce Uint256) {
	t.Nonce = nonce
}
