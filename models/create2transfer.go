package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

type Create2Transfer struct {
	TransactionBase
	ToStateID   *uint32
	ToPublicKey PublicKey
}

func (t *Create2Transfer) Type() txtype.TransactionType {
	return txtype.Create2Transfer
}

func (t *Create2Transfer) GetBase() *TransactionBase {
	return &t.TransactionBase
}

func (t *Create2Transfer) GetToStateID() *uint32 {
	return t.ToStateID
}

func (t *Create2Transfer) ToTransfer() *Transfer {
	panic("Create2Transfer cannot be cast to Transfer")
}

func (t *Create2Transfer) ToCreate2Transfer() *Create2Transfer {
	return t
}

func (t *Create2Transfer) ToMassMigration() *MassMigration {
	panic("Create2Transfer cannot be cast to MassMigration")
}

//nolint:gocritic
func (t Create2Transfer) Copy() GenericTransaction {
	return &t
}

//nolint:gocritic
func (t Create2Transfer) Clone() *Create2Transfer {
	return &t
}
