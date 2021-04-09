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
			err := XXX(storage, client, cfg)
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

func XXX(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	tx, txStorage, err := storage.BeginTransaction()
	if err != nil {
		return
	}
	defer tx.Rollback(&err)

	pendingTransactions, err := storage.GetPendingTransactions()
	if err != nil {
		return
	}

	commitments, err := commitmentsLoop(pendingTransactions, txStorage, cfg)
	if err != nil {
		return
	}

	err = SubmitBatch(commitments, txStorage, client, cfg)
	if err != nil {
		return err
	}
	
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func SubmitBatch(commitments []models.Commitment, storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	if len(commitments) < int(cfg.MinCommitmentsPerBatch) {
		return
	}

	batch, accountRoot, err := client.SubmitTransfersBatch(commitments)
	if err != nil {
		return
	}

	err = storage.AddBatch(batch)
	if err != nil {
		return
	}

	err = markCommitmentsAsIncluded(storage, commitments, &batch.Hash, accountRoot)
	if err != nil {
		return
	}

	log.Printf("Submitted %d commitment(s) on chain. Batch ID: %d. Batch Hash: %v", len(commitments), batch.ID.Uint64(), batch.Hash)

	return nil
}

func commitmentsLoop(
	pendingTransactions []models.Transaction,
	storage *st.Storage,
	cfg *config.RollupConfig,
) ([]models.Commitment, error) {
	stateTree := st.NewStateTree(storage)

	commitments := make([]models.Commitment, 0, 32)
	for {
		if len(commitments) >= int(cfg.MaxCommitmentsPerBatch) {
			break
		}

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return nil, err
		}

		includedTxs, err := ApplyTransactions(storage, pendingTransactions, cfg)
		if err != nil {
			return nil, err
		}

		oldPendingTxs := pendingTransactions
		pendingTransactions = make([]models.Transaction, 0)
		for i := range oldPendingTxs {
			tx := oldPendingTxs[i]
			included := false

			for i := range includedTxs {
				includedTx := includedTxs[i]
				if includedTx.Hash == tx.Hash {
					included = true
				}
			}

			if !included {
				pendingTransactions = append(pendingTransactions, tx)
			}
		}

		if len(includedTxs) < int(cfg.TxsPerCommitment) {
			err = stateTree.RevertTo(*initialStateRoot)
			if err != nil {
				return nil, err
			}
			break
		}

		log.Printf("Creating a commitment from %d transactions", len(includedTxs))
		commitment, err := createAndStoreCommitment(storage, includedTxs, cfg.FeeReceiverIndex)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)

		err = markTransactionsAsIncluded(storage, includedTxs, commitment.ID)
		if err != nil {
			return nil, err
		}
	}

	return commitments, nil
}

func markCommitmentsAsIncluded(storage *st.Storage, commitments []models.Commitment, batchHash, accountRoot *common.Hash) error {
	for i := range commitments {
		err := storage.MarkCommitmentAsIncluded(commitments[i].ID, batchHash, accountRoot)
		if err != nil {
			return err
		}
	}
	return nil
}
