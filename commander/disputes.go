package commander

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (c *Commander) handleBatchRollback(rollupCancel context.CancelFunc) (bool, error) {
	invalidBatchID, err := c.client.GetInvalidBatchID()
	if err != nil {
		return false, err
	}
	if invalidBatchID.IsZero() {
		c.invalidBatchID = *invalidBatchID
		return false, nil
	}

	return true, c.manageRemoteBatchRollback(invalidBatchID, rollupCancel)
}

func (c *Commander) manageRemoteBatchRollback(batchID *models.Uint256, rollupCancel context.CancelFunc) error {
	c.invalidBatchID = *batchID
	c.stopRollupLoop(rollupCancel)
	err := c.client.KeepRollingBack()
	if err != nil {
		return err
	}

	invalidBatch, err := c.storage.GetBatch(*batchID)
	if st.IsNotFoundError(err) {
		return nil
	}
	if err != nil {
		return err
	}

	return c.revertBatches(invalidBatch)
}
