package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Transfer struct {
	TransactionBase
	ToStateID uint32
}

type TransferForCommitment struct {
	TransactionBaseForCommitment
	ToStateID uint32
}

type TransferWithBatchDetails struct {
	Transfer
	BatchHash *common.Hash
	BatchTime *Timestamp
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
