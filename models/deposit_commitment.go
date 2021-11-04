package models

import (
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/zerohash"
	"github.com/ethereum/go-ethereum/common"
)

type DepositCommitment struct {
	CommitmentBase
	SubTreeID   Uint256
	SubTreeRoot common.Hash
	Deposits    []PendingDeposit
}

func (c *DepositCommitment) GetBodyHash() common.Hash {
	return zerohash.ZeroHash
}

func (c *DepositCommitment) LeafHash() common.Hash {
	return utils.HashTwo(c.PostStateRoot, c.GetBodyHash())
}
