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
	StateID  MerklePath
	DataHash common.Hash
	UserState
}

type StateUpdate struct {
	ID          uint64 `badgerhold:"key"`
	StateID     MerklePath
	CurrentHash common.Hash
	CurrentRoot common.Hash `badgerhold:"index"`
	PrevHash    common.Hash
	PrevRoot    common.Hash
}
