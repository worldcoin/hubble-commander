package models

import (
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentInclusionProof struct {
	StateRoot *common.Hash
	BodyRoot  *common.Hash
	Path      *MerklePath
	Witness   merkletree.Witness
}

type TransferCommitmentInclusionProof struct {
	StateRoot *common.Hash
	Transfer  *Transfer
	Path      *MerklePath
	Witness   merkletree.Witness
}

type StateMerkleProof struct {
	UserState *UserState
	Witness   merkletree.Witness
}
