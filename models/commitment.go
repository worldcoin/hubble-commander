package models

import "github.com/ethereum/go-ethereum/common"

type Commitment struct {
	LeafHash          common.Hash `db:"leaf_hash"`
	BodyHash          common.Hash `db:"body_hash"`
	AccountTreeRoot   common.Hash `db:"account_tree_root"`
	CombinedSignature Signature   `db:"combined_signature"`
	FeeReceiver       Uint256     `db:"fee_receiver"`
	Transactions      []byte
}
