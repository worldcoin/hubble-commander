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
	DataHash   common.Hash
}

type StateLeaf struct {
	StateID  uint32
	DataHash common.Hash
	UserState
}

type StateUpdate struct {
	ID            uint64 `badgerhold:"key"`
	CurrentRoot   common.Hash
	PrevRoot      common.Hash
	PrevStateLeaf StateLeaf
}
