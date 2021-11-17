package syncer

import (
	"context"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrBatchSubmissionFailed = errors.New("previous submit batch transaction failed")

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

func (c *Context) syncExistingBatch(remoteBatch eth.DecodedBatch, localBatch *models.Batch) error {
	batch := remoteBatch.GetBase()
	if batch.TransactionHash == localBatch.TransactionHash {
		err := c.batchCtx.UpdateExistingBatch(remoteBatch, *localBatch.PrevStateRoot)
		if err != nil {
			return err
		}

		log.Printf(
			"Synced new existing batch. Batch ID: %d. Batch Hash: %v",
			batch.ID.Uint64(),
			batch.Hash,
		)
	} else {
		txSender, err := c.getTransactionSender(batch.TransactionHash)
		if err != nil {
			return err
		}
		if *txSender != c.client.Blockchain.GetAccount().From {
			return NewInconsistentBatchError(localBatch)
		} else {
			// TODO remove the above check and this error once we use contracts with batchID verification:
			//  https://github.com/thehubbleproject/hubble-contracts/pull/601
			return ErrBatchSubmissionFailed
		}
	}
	return nil
}

func (c *Context) getTransactionSender(txHash common.Hash) (*common.Address, error) {
	tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, err
	}
	sender, err := types.LatestSignerForChainID(tx.ChainId()).Sender(tx)
	if err != nil {
		return nil, err
	}
	return &sender, nil
}
