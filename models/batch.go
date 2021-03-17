package models

import "github.com/ethereum/go-ethereum/common"

type Batch struct {
	Hash              common.Hash `db:"batch_hash"` // root of tree containing all commitments included in this batch
	ID                Uint256     `db:"batch_id"`
	FinalisationBlock Uint256
}
