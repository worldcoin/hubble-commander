package models

import "github.com/ethereum/go-ethereum/common"

type FlatStateLeaf struct {
	StateID    uint32
	DataHash   common.Hash
	PubKeyID   uint32  `badgerhold:"index"`
	TokenIndex Uint256 `badgerhold:"index"` // TODO: Consider removing or updating to the tuple of (Pubkey; tokenIdx)
	Balance    Uint256
	Nonce      Uint256
}

func NewFlatStateLeaf(leaf *StateLeaf) FlatStateLeaf {
	return FlatStateLeaf{
		StateID:    leaf.StateID,
		DataHash:   leaf.DataHash,
		PubKeyID:   leaf.PubKeyID,
		TokenIndex: leaf.TokenIndex,
		Balance:    leaf.Balance,
		Nonce:      leaf.Nonce,
	}
}

func (l *FlatStateLeaf) StateLeaf() *StateLeaf {
	return &StateLeaf{
		StateID:  l.StateID,
		DataHash: l.DataHash,
		UserState: UserState{
			PubKeyID:   l.PubKeyID,
			TokenIndex: l.TokenIndex,
			Balance:    l.Balance,
			Nonce:      l.Nonce,
		},
	}
}
