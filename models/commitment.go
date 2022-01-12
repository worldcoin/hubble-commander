package models

import "github.com/ethereum/go-ethereum/common"

type Commitment interface {
	GetCommitmentBase() CommitmentBase
	GetBodyHash() common.Hash
	GetPostStateRoot() common.Hash
	LeafHash() common.Hash
	ToTxCommitment() *TxCommitment
	ToMMCommitment() *MMCommitment
	ToDepositCommitment() *DepositCommitment
}

type CommitmentWithTxs interface {
	SetBodyHash(accountRoot common.Hash)
	CalcBodyHash(accountRoot common.Hash) *common.Hash
	ToTxCommitmentWithTxs() *TxCommitmentWithTxs
	ToMMCommitmentWithTxs() *MMCommitmentWithTxs
}
