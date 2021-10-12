package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type TransferCommitmentInclusionProof struct {
	StateRoot common.Hash
	Body      *TransferBody
	Path      *models.MerklePath
	Witness   models.Witness
}

type StateMerkleProof struct {
	models.StateMerkleProof
}

type PublicKeyProof struct {
	models.PublicKeyProof
}

type TransferBody struct {
	AccountRoot  common.Hash
	Signature    models.Signature
	FeeReceiver  uint32
	Transactions interface{}
}
