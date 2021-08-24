package models

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

type FlatStateLeaf struct {
	StateID  uint32
	DataHash common.Hash
	PubKeyID uint32  `badgerhold:"index"`
	TokenID  Uint256 // TODO: Consider adding a tuple index of (PubKeyID; TokenID)
	Balance  Uint256
	Nonce    Uint256
}

func MakeFlatStateLeaf(leaf *StateLeaf) FlatStateLeaf {
	return FlatStateLeaf{
		StateID:  leaf.StateID,
		DataHash: leaf.DataHash,
		PubKeyID: leaf.PubKeyID,
		TokenID:  leaf.TokenID,
		Balance:  leaf.Balance,
		Nonce:    leaf.Nonce,
	}
}

func (l *FlatStateLeaf) StateLeaf() *StateLeaf {
	return &StateLeaf{
		StateID:  l.StateID,
		DataHash: l.DataHash,
		UserState: UserState{
			PubKeyID: l.PubKeyID,
			TokenID:  l.TokenID,
			Balance:  l.Balance,
			Nonce:    l.Nonce,
		},
	}
}

func (l *FlatStateLeaf) Bytes() []byte {
	b := make([]byte, 136)
	binary.BigEndian.PutUint32(b[0:4], l.StateID)
	copy(b[4:36], l.DataHash[:])
	binary.BigEndian.PutUint32(b[36:40], l.PubKeyID)
	copy(b[40:72], l.TokenID.Bytes())
	copy(b[72:104], l.Balance.Bytes())
	copy(b[104:136], l.Nonce.Bytes())
	return b
}

func (l *FlatStateLeaf) SetBytes(data []byte) error {
	l.StateID = binary.BigEndian.Uint32(data[0:4])
	l.DataHash.SetBytes(data[4:36])
	l.PubKeyID = binary.BigEndian.Uint32(data[36:40])
	l.TokenID.SetBytes(data[40:72])
	l.Balance.SetBytes(data[72:104])
	l.Nonce.SetBytes(data[104:136])
	return nil
}
