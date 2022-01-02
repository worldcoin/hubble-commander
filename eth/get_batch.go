package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

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
		Committer:         meta.Committer,
	}, nil
}
