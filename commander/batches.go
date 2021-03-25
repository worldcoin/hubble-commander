package commander

import (
	"errors"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
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
		return
	}

	return tx.Commit()
}

func unsafeSubmitBatch(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	commitments, err := storage.GetPendingCommitments(uint64(cfg.MaxCommitmentsPerBatch))
	if err != nil {
		return err
	}
	if len(commitments) < int(cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	batch, _, err := client.SubmitTransfersBatch(commitments)
	if err != nil {
		return err
	}

	err = storage.AddBatch(batch)
	if err != nil {
		return err
	}

	err = markCommitmentsAsIncluded(storage, commitments, batch.Hash)
	if err != nil {
		return err
	}

	log.Printf("Sumbmited %d commitment(s) on chain. Batch ID: %d. Batch Hash: %v", len(commitments), batch.ID.Uint64(), batch.Hash)
	return nil
}

func markCommitmentsAsIncluded(storage *st.Storage, commitments []models.Commitment, batchHash common.Hash) error {
	for i := range commitments {
		err := storage.MarkCommitmentAsIncluded(commitments[i].ID, batchHash)
		if err != nil {
			return err
		}
	}
	return nil
}
