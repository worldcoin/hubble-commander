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
	DataHash   common.Hash `db:"data_hash"`
}

type StateLeaf struct {
	DataHash common.Hash `db:"data_hash"`
	UserState
}

type StateUpdate struct {
	ID          uint64
	StateID     MerklePath  `db:"state_id"`
	CurrentHash common.Hash `db:"current_hash"`
	CurrentRoot common.Hash `db:"current_root"`
	PrevHash    common.Hash `db:"prev_hash"`
	PrevRoot    common.Hash `db:"prev_root"`
}
