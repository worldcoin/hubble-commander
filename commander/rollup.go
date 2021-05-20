package commander

import (
	"errors"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func RollupEndlessLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	done := make(chan bool)
	return RollupLoop(storage, client, cfg, done)
}

func RollupLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, done <-chan bool) (err error) {
	ticker := time.NewTicker(cfg.BatchLoopInterval)
	defer ticker.Stop()

	currentBatchType := txtype.Transfer

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			err = SyncBatches(storage, client, cfg) // TODO disable in benchmarks
			if err != nil {
				return err
			}

			if currentBatchType == txtype.Transfer {
				err = createAndSubmitBatch(currentBatchType, storage, client, cfg)
				currentBatchType = txtype.Create2Transfer
			} else {
				err = createAndSubmitBatch(currentBatchType, storage, client, cfg)
				currentBatchType = txtype.Transfer
			}

			if err != nil {
				var e *RollupError
				if errors.As(err, &e) {
					continue
				}
				return err
			}
		}
	}
}

func createAndSubmitBatch(batchType txtype.TransactionType, storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	err = unsafeCreateAndSubmitBatch(batchType, txStorage, client, cfg)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func unsafeCreateAndSubmitBatch(
	batchType txtype.TransactionType,
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
) (err error) {
	var commitments []models.Commitment

	domain, err := storage.GetDomain(client.ChainState.ChainID)
	if err != nil {
		return err
	}

	if batchType == txtype.Transfer {
		commitments, err = buildTransferCommitments(storage, cfg, *domain)
	} else {
		commitments, err = buildCreate2TransfersCommitments(storage, client, cfg, *domain)
	}
	if err != nil {
		return err
	}

	err = submitBatch(batchType, commitments, storage, client, cfg)
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
