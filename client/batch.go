package client

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID              models.Uint256
	Type            batchtype.BatchType
	TransactionHash common.Hash
	PrevStateRoot   common.Hash
	Commitments     []PendingCommitment
}

func (b *Batch) ToDTO() dto.PendingBatch {
	commitments := make([]dto.PendingCommitment, 0, len(b.Commitments))
	for i := range b.Commitments {
		commitments = append(commitments, b.Commitments[i].ToDTO())
	}

	return dto.PendingBatch{
		ID:              b.ID,
		Type:            b.Type,
		TransactionHash: b.TransactionHash,
		PrevStateRoot:   b.PrevStateRoot,
		Commitments:     commitments,
	}
}
