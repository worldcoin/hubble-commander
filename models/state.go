package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type UserState struct {
	AccountIndex Uint256 `db:"account_index"`
	TokenIndex   Uint256 `db:"token_index"`
	Balance      Uint256
	Nonce        Uint256
}

type StateNode struct {
	MerklePath string      `db:"merkle_path"`
	DataHash   common.Hash `db:"data_hash"`
}

type StateLeaf struct {
	DataHash common.Hash `db:"data_hash"`
	UserState
}

type StateUpdate struct {
	ID          uint64
	MerklePath  string      `db:"merkle_path"`
	CurrentHash common.Hash `db:"current_hash"`
	CurrentRoot common.Hash `db:"current_root"`
	PrevHash    common.Hash `db:"prev_hash"`
	PrevRoot    common.Hash `db:"prev_root"`
}

func NewStateLeaf(state UserState) (*StateLeaf, error) {
	encodedState, err := EncodeUserState(state)
	if err != nil {
		return nil, err
	}
	dataHash := crypto.Keccak256Hash(encodedState)
	return &StateLeaf{
		DataHash:  dataHash,
		UserState: state,
	}, nil
}
