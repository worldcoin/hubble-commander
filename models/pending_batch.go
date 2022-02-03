package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type PendingBatch struct {
	ID              Uint256
	Type            batchtype.BatchType
	TransactionHash common.Hash
	Commitments     []PendingCommitment
}

type PendingCommitment struct {
	Commitment
	Transactions GenericTransactionArray
}
