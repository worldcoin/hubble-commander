package models

import (
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
)

type MMCommitment struct {
	CommitmentBase
	FeeReceiver       uint32
	CombinedSignature Signature
	BodyHash          *common.Hash
	Meta              *MassMigrationMeta
	WithdrawRoot      common.Hash
}

func (c *MMCommitment) GetBodyHash() common.Hash {
	return *c.BodyHash
}

func (c *MMCommitment) LeafHash() common.Hash {
	return utils.HashTwo(c.PostStateRoot, *c.BodyHash)
}

type MMCommitmentWithTxs struct {
	MMCommitment
	Transactions []byte
}

func (c *MMCommitmentWithTxs) SetBodyHash(accountRoot common.Hash) {
	c.BodyHash = calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *MMCommitmentWithTxs) CalcBodyHash(accountRoot common.Hash) *common.Hash {
	return calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *MMCommitmentWithTxs) ToTxCommitmentWithTxs() *TxCommitmentWithTxs {
	panic("Cannot cast MMCommitmentWithTxs to TxCommitmentWithTxs")
}

func (c *MMCommitmentWithTxs) ToMMCommitmentWithTxs() *MMCommitmentWithTxs {
	return c
}
