package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type StateNode struct {
	MerklePath string      `db:"merkle_path"`
	DataHash   common.Hash `db:"data_hash"`
}

type StateLeaf struct {
	DataHash     common.Hash `db:"data_hash"`
	AccountIndex Uint256     `db:"account_index"`
	TokenIndex   Uint256     `db:"token_index"`
	Balance      Uint256
	Nonce        Uint256
}

type StateUpdate struct {
	ID          uint64
	MerklePath  string      `db:"merkle_path"`
	CurrentHash common.Hash `db:"current_hash"`
	CurrentRoot common.Hash `db:"current_root"`
	PrevHash    common.Hash `db:"prev_hash"`
	PrevRoot    common.Hash `db:"prev_root"`
}
