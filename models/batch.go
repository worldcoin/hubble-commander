package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID                int32 `db:"batch_id"`
	Type              txtype.TransactionType
	TransactionHash   common.Hash  `db:"transaction_hash"`
	Hash              *common.Hash `db:"batch_hash"` // root of tree containing all commitments included in this batch
	Number            *Uint256     `db:"batch_number"`
	FinalisationBlock *uint32      `db:"finalisation_block"` // nolint:misspell
}

type BatchWithAccountRoot struct {
	Batch
	AccountTreeRoot *common.Hash `db:"account_tree_root"`
}
