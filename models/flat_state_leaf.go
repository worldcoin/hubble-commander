package models

import "github.com/ethereum/go-ethereum/common"

type FlatStateLeaf struct {
	DataHash   common.Hash
	PubKeyID   uint32  `badgerhold:"index"`
	TokenIndex Uint256 `badgerhold:"index"`
	Balance    Uint256
	Nonce      Uint256
}

func NewFlatStateLeaf(leaf *StateLeaf) FlatStateLeaf {
	return FlatStateLeaf{
		DataHash:   leaf.DataHash,
		PubKeyID:   leaf.PubKeyID,
		TokenIndex: leaf.TokenIndex,
		Balance:    leaf.Balance,
		Nonce:      leaf.Nonce,
	}
}

func (l *FlatStateLeaf) StateLeaf() *StateLeaf {
	return &StateLeaf{
		DataHash: l.DataHash,
		UserState: UserState{
			PubKeyID:   l.PubKeyID,
			TokenIndex: l.TokenIndex,
			Balance:    l.Balance,
			Nonce:      l.Nonce,
		},
	}
}
