package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type UserState struct {
	PubkeyID   uint32  `db:"pubkey_id"`
	TokenIndex Uint256 `db:"token_index"`
	Balance    Uint256
	Nonce      Uint256
}

type StateNode struct {
	StateID  MerklePath  `db:"state_id"`
	DataHash common.Hash `db:"data_hash"`
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
