package commander

import (
	"errors"
	"time"

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
	tx, txStorage, err := storage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	var commitments []models.Commitment

	if batchType == txtype.Transfer {
		pendingTransfers, err := storage.GetPendingTransfers()
		if err != nil {
			return err
		}

		commitments, err = createTransferCommitments(pendingTransfers, txStorage, cfg)
		if err != nil {
			return err
		}
	} else {
		pendingTransfers, err := storage.GetPendingCreate2Transfers()
		if err != nil {
			return err
		}

		commitments, err = createCreate2TransferCommitments(pendingTransfers, txStorage, cfg)
		if err != nil {
			return err
		}
	}

	err = submitBatch(batchType, commitments, txStorage, client, cfg)
	if err != nil {
		return err
	}

	return tx.Commit()
}
