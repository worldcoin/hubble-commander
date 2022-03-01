package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID                Uint256
	Type              batchtype.BatchType
	TransactionHash   common.Hash
	Hash              *common.Hash // root of tree containing all commitments included in this batch
	FinalisationBlock *uint32
	AccountTreeRoot   *common.Hash
	PrevStateRoot     common.Hash
	MinedTime         *Timestamp
}
