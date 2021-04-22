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

type BatchWithCommitments struct {
	Batch
	Commitments []Commitment
}
