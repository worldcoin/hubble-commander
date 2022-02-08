package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type ContractBatch struct {
	ID   models.Uint256
	Hash common.Hash
	models.BatchMeta
}

func (cb *ContractBatch) ToModelBatch() *models.Batch {
	return &models.Batch{
		ID:                cb.ID,
		Type:              cb.BatchType,
		Hash:              &cb.Hash,
		FinalisationBlock: &cb.FinaliseOn,
	}
}

// TODO Replace usages of GetBatch with GetContractBatch

func (c *Client) GetBatch(batchID *models.Uint256) (*models.Batch, error) {
	batch, err := c.Rollup.GetBatch(nil, batchID.ToBig())
	if err != nil {
		return nil, err
	}
	meta := encoder.DecodeMeta(batch.Meta)
	hash := common.BytesToHash(batch.CommitmentRoot[:])
	return &models.Batch{
		ID:                *batchID,
		Hash:              &hash,
		Type:              meta.BatchType,
		FinalisationBlock: &meta.FinaliseOn,
	}, nil
}

func (c *Client) GetContractBatch(batchID *models.Uint256) (*ContractBatch, error) {
	batch, err := c.Rollup.GetBatch(nil, batchID.ToBig())
	if err != nil {
		return nil, err
	}
	meta := encoder.DecodeMeta(batch.Meta)
	hash := common.BytesToHash(batch.CommitmentRoot[:])
	return &ContractBatch{
		ID:        *batchID,
		Hash:      hash,
		BatchMeta: meta,
	}, nil
}
