package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) GetBatch(batchNumber *models.Uint256) (*models.Batch, error) {
	batch, err := c.Rollup.GetBatch(nil, batchNumber.ToBig())
	if err != nil {
		return nil, err
	}
	meta := encoder.DecodeMeta(batch.Meta)
	hash := common.BytesToHash(batch.CommitmentRoot[:])
	return &models.Batch{
		Hash:              &hash,
		Type:              meta.BatchType,
		Number:            batchNumber,
		FinalisationBlock: &meta.FinaliseOn,
	}, nil
}
