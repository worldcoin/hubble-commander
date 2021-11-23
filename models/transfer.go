package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Transfer struct {
	TransactionBase
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

func (t *Transfer) GetSpokeID() *Uint256 {
	panic("Transfer does not contain a SpokeID field")
}

func (t *Transfer) ToTransfer() *Transfer {
	return t
}

func (t *Transfer) ToCreate2Transfer() *Create2Transfer {
	panic("Transfer cannot be cast to Create2Transfer")
}

func (t *Transfer) ToMassMigration() *MassMigration {
	panic("Transfer cannot be cast to MassMigration")
}

// nolint:gocritic
func (t Transfer) Copy() GenericTransaction {
	return &t
}
