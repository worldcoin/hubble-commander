package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type Commitment struct {
	ID                 int32
	LeafHash           common.Hash    //commitment.LeafHash()
	TokenID            models.Uint256 //from feeReceiverStateID -> getStateLeafByPath
	FeeReceiverStateID uint32
	CombinedSignature  models.Signature
	PostStateRoot      common.Hash
}
