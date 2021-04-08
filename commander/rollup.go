package commander

import (
	"errors"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
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
			err := NAME_ME_BETTER(storage, client, cfg)
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

func NAME_ME_BETTER(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	tx, txStorage, err := storage.BeginTransaction()
	if err != nil {
		return
	}
	defer tx.Rollback(&err)

	pendingTransactions, err := storage.GetPendingTransactions()
	if err != nil {
		return
	}

	stateTree := st.NewStateTree(txStorage)

	commitments := make([]models.Commitment, 0, 32)
	for {
		if len(commitments) >= int(cfg.MaxCommitmentsPerBatch) {
			break
		}

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return err
		}

		includedTxs, _, err := ApplyTransactions(txStorage, pendingTransactions, cfg)
		if err != nil {
			return err
		}

		oldPendingTxs := pendingTransactions
		pendingTransactions = make([]models.Transaction, 0)
		for _, tx := range oldPendingTxs {
			included := false
			for _, inc := range includedTxs {
				if inc.Hash == tx.Hash {
					included = true
				}
			}

			if !included {
				pendingTransactions = append(pendingTransactions, tx)
			}
		}

		if len(includedTxs) < int(cfg.TxsPerCommitment) {
			stateTree.RevertTo(*initialStateRoot)
			break
		}

		log.Printf("Creating a commitment from %d transactions", len(includedTxs))
		commitment, err := createAndStoreCommitment(txStorage, includedTxs, cfg.FeeReceiverIndex)
		if err != nil {
			return err
		}

		commitments = append(commitments, *commitment)

		err = markTransactionsAsIncluded(txStorage, includedTxs, commitment.ID)
		if err != nil {
			return err
		}
	}

	if len(commitments) < int(cfg.MinCommitmentsPerBatch) {
		return
	}

	batch, accountRoot, err := client.SubmitTransfersBatch(commitments)
	if err != nil {
		return
	}

	err = storage.AddBatch(batch)
	if err != nil {
		return err
	}

	err = markCommitmentsAsIncluded(txStorage, commitments, &batch.Hash, accountRoot)
	if err != nil {
		return err
	}

	log.Printf("Submitted %d commitment(s) on chain. Batch ID: %d. Batch Hash: %v", len(commitments), batch.ID.Uint64(), batch.Hash)

	tx.Commit()

	return nil
}
