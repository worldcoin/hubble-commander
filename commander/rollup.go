package commander

import (
	"errors"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (c *Commander) rollupLoop() (err error) {
	ticker := time.NewTicker(c.cfg.Rollup.BatchLoopInterval)
	defer ticker.Stop()

	currentBatchType := txtype.Transfer

	for {
		select {
		case <-c.stopChannel:
			return nil
		case <-ticker.C:
			if c.cfg.Rollup.SyncBatches {
				err = SyncBatches(c.storage, c.client, &c.cfg.Rollup)
				if err != nil {
					return err
				}
			}

			var transactionExecutor *transactionExecutor
			transactionExecutor, err = newTransactionExecutor(c.storage, c.client, &c.cfg.Rollup)
			if err != nil {
				return err
			}

			if currentBatchType == txtype.Transfer {
				err = transactionExecutor.CreateAndSubmitBatch(currentBatchType)
				currentBatchType = txtype.Create2Transfer
			} else {
				err = transactionExecutor.CreateAndSubmitBatch(currentBatchType)
				currentBatchType = txtype.Transfer
			}

			if err != nil {
				transactionExecutor.Rollback(&err)

				var e *RollupError
				if errors.As(err, &e) {
					continue
				}
				return err
			}

			err = transactionExecutor.Commit()
			if err != nil {
				return err
			}
		}
	}
}

// transactionExecutor executes transactions & syncs batches. Manages a database transaction.
type transactionExecutor struct {
	cfg     *config.RollupConfig
	storage *st.Storage
	tx      *db.TxController
	client  *eth.Client
}

func newTransactionExecutor(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (*transactionExecutor, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	return &transactionExecutor{
		cfg:     cfg,
		storage: txStorage,
		tx:      tx,
		client:  client,
	}, nil
}

func (t *transactionExecutor) Commit() error {
	return t.tx.Commit()
}

func (t *transactionExecutor) Rollback(cause *error) {
	t.tx.Rollback(cause)
}

func (t *transactionExecutor) CreateAndSubmitBatch(batchType txtype.TransactionType) (err error) {
	var commitments []models.Commitment

	domain, err := t.storage.GetDomain(t.client.ChainState.ChainID)
	if err != nil {
		return err
	}

	if batchType == txtype.Transfer {
		commitments, err = buildTransferCommitments(t.storage, t.cfg, *domain)
	} else {
		commitments, err = buildCreate2TransfersCommitments(t.storage, t.client, t.cfg, *domain)
	}
	if err != nil {
		return err
	}

	err = submitBatch(batchType, commitments, t.storage, t.client, t.cfg)
	if err != nil {
		return err
	}
	return nil
}

func buildTransferCommitments(storage *st.Storage, cfg *config.RollupConfig, domain bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := storage.GetPendingTransfers()
	if err != nil {
		return nil, err
	}
	return createTransferCommitments(pendingTransfers, storage, cfg, domain)
}

func buildCreate2TransfersCommitments(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	domain bls.Domain,
) ([]models.Commitment, error) {
	pendingTransfers, err := storage.GetPendingCreate2Transfers()
	if err != nil {
		return nil, err
	}
	return createCreate2TransferCommitments(pendingTransfers, storage, client, cfg, domain)
}
