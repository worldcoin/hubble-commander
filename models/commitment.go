package models

import "github.com/ethereum/go-ethereum/common"

type Commitment struct {
	LeafHash          common.Hash `db:"leaf_hash"`
	PostStateRoot     common.Hash `db:"post_state_root"`
	BodyHash          common.Hash `db:"body_hash"`
	AccountTreeRoot   common.Hash `db:"account_tree_root"`
	CombinedSignature Signature   `db:"combined_signature"`
	FeeReceiver       uint32      `db:"fee_receiver"`
	Transactions      []byte
}
