package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

type MassMigration struct {
	TransactionBase
	SpokeID uint32
}

func (m *MassMigration) Type() txtype.TransactionType {
	return txtype.MassMigration
}

func (m *MassMigration) GetBase() *TransactionBase {
	return &m.TransactionBase
}

func (m *MassMigration) GetToStateID() *uint32 {
	panic("MassMigration does not contain a ToStateID field")
}

func (m *MassMigration) ToTransfer() *Transfer {
	panic("MassMigration cannot be cast to Transfer")
}

func (m *MassMigration) ToCreate2Transfer() *Create2Transfer {
	panic("MassMigration cannot be cast to Create2Transfer")
}

func (m *MassMigration) ToMassMigration() *MassMigration {
	return m
}

//nolint:gocritic
func (m MassMigration) Copy() GenericTransaction {
	return &m
}
