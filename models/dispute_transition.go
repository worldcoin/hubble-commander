package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentInclusionProof struct {
	StateRoot common.Hash
	BodyRoot  common.Hash
	Path      *MerklePath
	Witness   Witness
}

type TransferCommitmentInclusionProof struct {
	StateRoot common.Hash
	Body      *TransferBody
	Path      *MerklePath
	Witness   Witness
}

type TransferBody struct {
	AccountRoot  common.Hash
	Signature    Signature
	FeeReceiver  uint32
	Transactions []byte
}

type MMCommitmentInclusionProof struct {
	StateRoot common.Hash
	Body      *MMBody
	Path      *MerklePath
	Witness   Witness
}

type MMBody struct {
	AccountRoot  common.Hash
	Signature    Signature
	Meta         MassMigrationMeta
	WithdrawRoot common.Hash
	Transactions []byte
}

type StateMerkleProof struct {
	UserState *UserState
	Witness   Witness
}
