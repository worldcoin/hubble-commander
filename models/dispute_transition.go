package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentInclusionProof struct {
	StateRoot common.Hash
	BodyRoot  common.Hash
	Path      *MerklePath
	Witnesses Witnesses
}

type TransferCommitmentInclusionProof struct {
	StateRoot common.Hash
	Body      *TransferBody
	Path      *MerklePath
	Witnesses Witnesses
}

type TransferBody struct {
	AccountRoot  common.Hash
	Signature    Signature
	FeeReceiver  uint32
	Transactions []byte
}

type StateMerkleProof struct {
	UserState *UserState
	Witnesses Witnesses
}
