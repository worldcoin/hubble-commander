package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentInclusionProofBase struct {
	StateRoot common.Hash
	Path      *MerklePath
	Witness   models.Witness
}

type CommitmentInclusionProof struct {
	CommitmentInclusionProofBase
	Body *CommitmentProofBody
}

type MassMigrationCommitmentProof struct {
	CommitmentInclusionProofBase
	Body *MassMigrationBody
}

type MassMigrationBody struct {
	AccountRoot  common.Hash
	Signature    models.Signature
	Meta         *MassMigrationMeta
	WithdrawRoot common.Hash
	Transactions []byte
}

type StateMerkleProof struct {
	UserState *UserState
	Witness   models.Witness
}

type WithdrawProof struct {
	UserState *UserState
	Path      MerklePath
	Witness   models.Witness
	Root      common.Hash
}

type PublicKeyProof struct {
	PublicKey *models.PublicKey
	Witness   models.Witness
}

type CommitmentProofBody struct {
	AccountRoot  common.Hash
	Signature    models.Signature
	FeeReceiver  uint32
	Transactions interface{}
}
