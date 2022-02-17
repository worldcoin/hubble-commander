package models

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionBase struct {
	Hash           common.Hash
	TxType         txtype.TransactionType
	FromStateID    uint32
	Amount         Uint256
	Fee            Uint256
	Nonce          Uint256
	Signature      Signature
	ReceiveTime    *Timestamp
	CommitmentSlot *CommitmentSlot
	ErrorMessage   *string
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

func (t *TransactionBase) SetReceiveTime() {
	t.ReceiveTime = NewTimestamp(time.Now().UTC())
}
