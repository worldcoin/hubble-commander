package syncer

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
)

func (c *Context) SyncBatch(remoteBatch eth.DecodedBatch) error {
	localBatch, err := c.storage.GetBatch(remoteBatch.GetID())
	if err != nil && !st.IsNotFoundError(err) {
		return err
	}

	if st.IsNotFoundError(err) {
		return c.syncNewBatch(remoteBatch)
	} else {
		return c.syncExistingBatch(remoteBatch, localBatch)
	}
}

func (c *Context) syncNewBatch(remoteBatch eth.DecodedBatch) error {
	logSyncStart(remoteBatch)

	root, err := c.storage.StateTree.Root()
	if err != nil {
		return err
	}

	err = c.storage.AddBatch(remoteBatch.ToBatch(*root))
	if err != nil {
		return err
	}

	err = c.batchCtx.SyncCommitments(remoteBatch)
	if err != nil {
		return err
	}

	logSyncSuccess(remoteBatch)
	return nil
}

func logSyncStart(batch eth.DecodedBatch) {
	log.Debugf("Syncing new %s batch #%s with %d commitment(s) from chain",
		batch.GetBase().Type.String(),
		batch.GetBase().ID.String(),
		batch.GetCommitmentsLength(),
	)
}

func logSyncSuccess(batch eth.DecodedBatch) {
	log.Printf("Synced new %s batch #%s with %d commitment(s) from chain",
		batch.GetBase().Type.String(),
		batch.GetBase().ID.String(),
		batch.GetCommitmentsLength(),
	)
}

func (c *Context) syncExistingBatch(remoteDecodedBatch eth.DecodedBatch, localBatch *models.Batch) error {
	remoteBatch := remoteDecodedBatch.GetBase()
	if remoteBatch.TransactionHash == localBatch.TransactionHash {
		err := c.batchCtx.UpdateExistingBatch(remoteDecodedBatch, *localBatch.PrevStateRoot)
		if err != nil {
			return err
		}
		log.Printf("Synced new existing batch. Batch ID: %s. Batch Hash: %s", remoteBatch.ID.String(), remoteBatch.Hash.String())
	} else {
		// This can happen when proposer slot ends before batch submission transaction gets mined
		log.Errorf("Local batch with ID=%s was found inconsistent with remote batch", localBatch.ID.String())
		return NewInconsistentBatchError(localBatch)
	}

	err := c.storage.AddPendingStakeWithdrawal(&models.PendingStakeWithdrawal{
		BatchID:           localBatch.ID,
		FinalisationBlock: remoteBatch.FinalisationBlock,
	})
	if err != nil {
		return err
	}
	return nil
}
