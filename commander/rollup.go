package commander

import (
	"errors"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
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
				err = createAndSubmitTransferBatch(storage, client, cfg)
				currentBatchType = txtype.Create2Transfer
			} else {
				err = createAndSubmitCreate2TransferBatch(storage, client, cfg)
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

func createAndSubmitTransferBatch(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	tx, txStorage, err := storage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	pendingTransfers, err := storage.GetPendingTransfers()
	if err != nil {
		return err
	}

	commitments, err := createTransferCommitments(pendingTransfers, txStorage, cfg)
	if err != nil {
		return err
	}

	err = submitTransferBatch(commitments, txStorage, client, cfg)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func createAndSubmitCreate2TransferBatch(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	tx, txStorage, err := storage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	pendingTransfers, err := storage.GetPendingCreate2Transfers()
	if err != nil {
		return err
	}

	commitments, err := createCreate2TransferCommitments(pendingTransfers, txStorage, cfg)
	if err != nil {
		return err
	}

	err = submitCreate2TransferBatch(commitments, txStorage, client, cfg)
	if err != nil {
		return err
	}

	return tx.Commit()
}
