package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID                Uint256 `db:"batch_id"`
	Type              txtype.TransactionType
	TransactionHash   common.Hash  `db:"transaction_hash"`
	Hash              *common.Hash `db:"batch_hash"`         // root of tree containing all commitments included in this batch
	FinalisationBlock *uint32      `db:"finalisation_block"` // nolint:misspell
	AccountTreeRoot   *common.Hash `db:"account_tree_root"`
	PrevStateRoot     *common.Hash `db:"prev_state_root"`
	SubmissionTime    *Timestamp   `db:"submission_time"`
}
