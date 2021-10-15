package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type DepositCommitment struct {
	ID            CommitmentID
	Type          batchtype.BatchType
	PostStateRoot common.Hash
	SubTreeID     Uint256
	SubTreeRoot   common.Hash
	Deposits      []PendingDeposit
}
