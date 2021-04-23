package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID                models.Uint256
	Hash              common.Hash
	Type              txtype.TransactionType
	AccountTreeRoot   *common.Hash
	FinalisationBlock uint32
}

type BatchWithCommitments struct {
	Batch
	Commitments []Commitment
}
