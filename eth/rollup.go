package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) GetBatch(batchID *models.Uint256) (*models.Batch, error) {
	batch, err := c.Rollup.GetBatch(nil, &batchID.Int)
	if err != nil {
		return nil, err
	}
	meta := encoder.DecodeMeta(batch.Meta)
	return &models.Batch{
		Hash:              common.BytesToHash(batch.CommitmentRoot[:]),
		Type:              meta.BatchType,
		ID:                *batchID,
		FinalisationBlock: meta.FinaliseOn,
	}, nil
}
