package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type UserState struct {
	PubKeyID   uint32
	TokenIndex Uint256
	Balance    Uint256
	Nonce      Uint256
}

type StateNode struct {
	MerklePath MerklePath
	DataHash   common.Hash `badgerhold:"index"`
}

type StateLeaf struct {
	StateID  uint32
	DataHash common.Hash
	UserState
}

type StateUpdate struct {
	ID            uint64 `badgerhold:"key"`
	StateID       uint32
	CurrentRoot   common.Hash `badgerhold:"index"`
	PrevRoot      common.Hash
	PrevStateLeaf StateLeaf
}
