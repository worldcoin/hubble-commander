package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentInclusionProof struct {
	StateRoot common.Hash
	Body      *CommitmentProofBody
	Path      *models.MerklePath
	Witness   models.Witness
}

type StateMerkleProof struct {
	models.StateMerkleProof
}

type PublicKeyProof struct {
	models.PublicKeyProof
}

type CommitmentProofBody struct {
	AccountRoot  common.Hash
	Signature    models.Signature
	FeeReceiver  uint32
	Transactions interface{}
}
