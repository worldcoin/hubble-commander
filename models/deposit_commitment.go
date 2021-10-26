package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type DepositCommitment struct {
	CommitmentBase
	SubTreeID   Uint256
	SubTreeRoot common.Hash
	Deposits    []PendingDeposit
}
