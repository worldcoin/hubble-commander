package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentInclusionProof struct {
	StateRoot common.Hash
	Body      *CommitmentProofBody
	Path      *MerklePath
	Witness   models.Witness
}

type StateMerkleProof struct {
	UserState *UserState
	Witness   models.Witness
}

type WithdrawProof struct {
	models.WithdrawProof
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
