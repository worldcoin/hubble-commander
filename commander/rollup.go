package commander

import (
	"errors"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func RollupEndlessLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	done := make(chan bool)
	return RollupLoop(storage, client, cfg, done)
}

func RollupLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, done <-chan bool) error {
	ticker := time.NewTicker(cfg.BatchLoopInterval)

	for {
		select {
		case <-done:
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := commitTransactions(storage, client, cfg)
			if err != nil {
				var e *BatchError
				if errors.As(err, &e) {
					log.Println(e.Error())
					continue
				}
				return err
			}
		}
	}
}

func commitTransactions(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	tx, txStorage, err := storage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	pendingTransactions, err := storage.GetPendingTransactions()
	if err != nil {
		return err
	}

	commitments, err := createCommitments(pendingTransactions, txStorage, cfg)
	if err != nil {
		return err
	}

	err = submitBatch(commitments, txStorage, client, cfg)
	if err != nil {
		// ! Temporary solutions - handle it better when other types of batches are added
		if err != ErrNotEnoughCommitments {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
