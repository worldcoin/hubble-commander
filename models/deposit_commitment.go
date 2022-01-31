package models

import (
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/ethereum/go-ethereum/common"
)

type DepositCommitment struct {
	CommitmentBase
	SubtreeID   Uint256
	SubtreeRoot common.Hash
	Deposits    []PendingDeposit
}

func (c *DepositCommitment) GetCommitmentBase() *CommitmentBase {
	return &c.CommitmentBase
}

func (c *DepositCommitment) SetBodyHash(_ *common.Hash) {
	// NOOP
}

func (c *DepositCommitment) GetBodyHash() *common.Hash {
	return &consts.ZeroHash
}

func (c *DepositCommitment) LeafHash() common.Hash {
	return utils.HashTwo(c.PostStateRoot, *c.GetBodyHash())
}

func (c *DepositCommitment) ToTxCommitment() *TxCommitment {
	panic("cannot cast DepositCommitment to TxCommitment")
}

func (c *DepositCommitment) ToMMCommitment() *MMCommitment {
	panic("cannot cast DepositCommitment to MMCommitment")
}

func (c *DepositCommitment) ToDepositCommitment() *DepositCommitment {
	return c
}
