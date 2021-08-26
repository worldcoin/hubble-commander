package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionBase struct {
	Hash         common.Hash
	TxType       txtype.TransactionType
	FromStateID  uint32
	Amount       Uint256
	Fee          Uint256
	Nonce        Uint256
	Signature    Signature
	ReceiveTime  *Timestamp
	CommitmentID *CommitmentID
	ErrorMessage *string
}

type TransactionBaseForCommitment struct {
	Hash        common.Hash `db:"tx_hash"`
	FromStateID uint32      `db:"from_state_id"`
	Amount      Uint256
	Fee         Uint256
	Nonce       Uint256
	Signature   Signature
	ReceiveTime *Timestamp `db:"receive_time"`
}

func (t *TransactionBase) GetFromStateID() uint32 {
	return t.FromStateID
}

func (t *TransactionBase) GetAmount() Uint256 {
	return t.Amount
}

func (t *TransactionBase) GetFee() Uint256 {
	return t.Fee
}

func (t *TransactionBase) GetNonce() Uint256 {
	return t.Nonce
}

func (t *TransactionBase) SetNonce(nonce Uint256) {
	t.Nonce = nonce
}

func (t *TransactionBase) GetSignature() Signature {
	return t.Signature
}
