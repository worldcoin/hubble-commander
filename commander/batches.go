package commander

import (
	"errors"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

var (
	ErrNotEnoughCommitments = NewBatchError("not enough commitments")
)

func BatchesEndlessLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	done := make(chan bool)
	return BatchesLoop(storage, client, cfg, done)
}

func BatchesLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, done <-chan bool) error {
	ticker := time.NewTicker(cfg.BatchLoopInterval)

	for {
		select {
		case <-done:
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := SubmitBatch(storage, client, cfg)
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

func SubmitBatch(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	tx, txStorage, err := storage.BeginTransaction()
	if err != nil {
		return
	}
	defer tx.Rollback(&err)

	err = unsafeSubmitBatch(txStorage, client, cfg)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// nolint:unparam
func unsafeSubmitBatch(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	commitments, err := storage.GetPendingCommitments(uint64(cfg.MaxCommitmentsPerBatch))
	if err != nil {
		return err
	}
	if len(commitments) < int(cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}
	return nil
}
