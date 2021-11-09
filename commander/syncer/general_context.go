package syncer

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

type GeneralContext struct {
	storage *st.Storage
	tx      *db.TxController
	client  *eth.Client
	syncCtx SyncContext
}

func NewGeneralContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) *GeneralContext {
	tx, txStorage := storage.BeginTransaction(st.TxOptions{})

	switch batchType {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return newGeneralContext(txStorage, tx, client, newContext(txStorage, client, cfg, batchType))
	case batchtype.Deposit:
		return newGeneralContext(txStorage, tx, client, newDepositsContext(txStorage, client, cfg))
	case batchtype.Genesis, batchtype.MassMigration:
		panic("invalid batch type")
	}
	return nil
}

func NewTestGeneralContext(
	storage *st.Storage,
	tx *db.TxController,
	client *eth.Client,
	syncCtx SyncContext,
) *GeneralContext {
	return newGeneralContext(storage, tx, client, syncCtx)
}

func (c *GeneralContext) Commit() error {
	return c.tx.Commit()
}

// nolint:gocritic
func (c *GeneralContext) Rollback(cause *error) {
	c.tx.Rollback(cause)
}

func newGeneralContext(
	storage *st.Storage,
	tx *db.TxController,
	client *eth.Client,
	syncCtx SyncContext,
) *GeneralContext {
	return &GeneralContext{
		storage: storage,
		tx:      tx,
		client:  client,
		syncCtx: syncCtx,
	}
}

func (c *GeneralContext) SyncBatch(remoteBatch eth.DecodedBatch) error {
	localBatch, err := c.storage.GetBatch(remoteBatch.GetBatch().ID)
	if err != nil && !st.IsNotFoundError(err) {
		return err
	}

	if st.IsNotFoundError(err) {
		return c.syncCtx.SyncNewBatch(remoteBatch)
	} else {
		return c.syncExistingBatch(remoteBatch, localBatch)
	}
}

func (c *GeneralContext) syncExistingBatch(remoteBatch eth.DecodedBatch, localBatch *models.Batch) error {
	batch := remoteBatch.GetBatch()
	if batch.TransactionHash == localBatch.TransactionHash {
		err := c.syncCtx.UpdateExistingBatch(remoteBatch)
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

func (c *GeneralContext) getTransactionSender(txHash common.Hash) (*common.Address, error) {
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
