package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type UserState struct {
	PubKeyID   uint32  `db:"pub_key_id"`
	TokenIndex Uint256 `db:"token_index"`
	Balance    Uint256
	Nonce      Uint256
}

type StateNode struct {
	MerklePath MerklePath  `db:"merkle_path"`
	DataHash   common.Hash `db:"data_hash" badgerhold:"index"`
}

type StateLeaf struct {
	DataHash common.Hash `db:"data_hash"`
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
