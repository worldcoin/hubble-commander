package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	Hash              common.Hash `db:"batch_hash"` // root of tree containing all commitments included in this batch
	Type              txtype.TransactionType
	ID                Uint256 `db:"batch_id"`
	FinalisationBlock uint32  `db:"finalisation_block"` // nolint:misspell
}

type BatchWithAccountRoot struct {
	ID                Uint256
	Hash              common.Hash
	Type              txtype.TransactionType
	AccountTreeRoot   *common.Hash `db:"account_tree_root"`
	FinalisationBlock uint32       `db:"finalisation_block"`
}
