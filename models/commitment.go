package models

import "github.com/ethereum/go-ethereum/common"

type Commitment struct {
	LeafHash          common.Hash `db:"leaf_hash"`
	BodyHash          common.Hash `db:"body_hash"`
	AccountTreeRoot   common.Hash `db:"account_tree_root"`
	CombinedSignature []byte      `db:"combined_signature"` // TODO: Define type for signatures
	FeeReceiver       Uint256     `db:"fee_receiver"`
	Transactions      []byte
}
