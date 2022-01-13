package models

import "github.com/ethereum/go-ethereum/common"

type Commitment interface {
	GetCommitmentBase() *CommitmentBase
	SetBodyHash(bodyHash *common.Hash)
	GetBodyHash() common.Hash
	GetPostStateRoot() common.Hash
	LeafHash() common.Hash
	ToTxCommitment() *TxCommitment
	ToMMCommitment() *MMCommitment
	ToDepositCommitment() *DepositCommitment
}

type CommitmentWithTxs interface {
	CalcBodyHash(accountRoot common.Hash) *common.Hash
	CalcAndSetBodyHash(accountRoot common.Hash)
	ToTxCommitmentWithTxs() *TxCommitmentWithTxs
	ToMMCommitmentWithTxs() *MMCommitmentWithTxs
}
