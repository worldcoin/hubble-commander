package commander

import (
	"errors"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func RollupEndlessLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	done := make(chan bool)
	return RollupLoop(storage, client, cfg, done)
}

func RollupLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, done <-chan bool) (err error) {
	ticker := time.NewTicker(cfg.BatchLoopInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			err = createAndSubmitBatch(storage, client, cfg)
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

func createAndSubmitBatch(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
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
