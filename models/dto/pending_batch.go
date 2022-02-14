package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type PendingBatch struct {
	ID              models.Uint256
	Type            batchtype.BatchType
	TransactionHash common.Hash
	PrevStateRoot   common.Hash
	Commitments     []PendingCommitment
}

type PendingCommitment struct {
	models.Commitment

	// We're using type from models, because commander is only consumer of API returning that type
	Transactions models.GenericTransactionArray
}
