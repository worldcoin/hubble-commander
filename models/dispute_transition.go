package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentInclusionProofBase struct {
	StateRoot common.Hash
	Path      *MerklePath
	Witness   Witness
}

type CommitmentInclusionProof struct {
	CommitmentInclusionProofBase
	BodyRoot common.Hash
}

type TransferCommitmentInclusionProof struct {
	CommitmentInclusionProofBase
	Body *TransferBody
}

type TransferBody struct {
	AccountRoot  common.Hash
	Signature    Signature
	FeeReceiver  uint32
	Transactions []byte
}

type MMCommitmentInclusionProof struct {
	CommitmentInclusionProofBase
	Body *MMBody
}

type MMBody struct {
	AccountRoot  common.Hash
	Signature    Signature
	Meta         *MassMigrationMeta
	WithdrawRoot common.Hash
	Transactions []byte
}

type StateMerkleProof struct {
	UserState *UserState
	Witness   Witness
}

type WithdrawProof struct {
	UserState *UserState
	Path      MerklePath
	Witness   Witness
	Root      common.Hash
}
