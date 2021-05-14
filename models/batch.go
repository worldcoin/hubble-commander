package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	Hash              common.Hash `db:"batch_hash"` // root of tree containing all commitments included in this batch
	Type              txtype.TransactionType
	ID                Uint256 `db:"batch_id"`           // TODO consider adding an index
	FinalisationBlock uint32  `db:"finalisation_block"` // nolint:misspell
}

type BatchWithSubmissionBlock struct {
	Batch
	SubmissionBlock uint32
}

type BatchWithAccountRoot struct {
	BatchWithSubmissionBlock
	AccountTreeRoot *common.Hash `db:"account_tree_root"`
}
