package commander

import (
	"fmt"
	"log"
	"sync"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var (
	ErrFraudulentTransfer    = errors.New("fraudulent transfer encountered when syncing")
	ErrTransfersNotApplied   = errors.New("could not apply all transfers from synced batch")
	ErrBatchSubmissionFailed = errors.New("previous submit batch transaction failed")
)

func (t *transactionExecutor) SyncBatches(stateMutex *sync.Mutex, startBlock, endBlock uint64) error {
	latestBatchNumber, err := getLatestBatchNumber(t.storage)
	if err != nil {
		return err
	}

	newBatches, err := t.client.GetBatches(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	})
	if err != nil {
		return err
	}

	for i := range newBatches {
		batch := &newBatches[i]
		if batch.Number.Cmp(latestBatchNumber) <= 0 {
			continue
		}

		localBatch, err := t.storage.GetBatchByNumber(batch.Number)
		if err != nil && !st.IsNotFoundError(err) {
			return err
		}

		if st.IsNotFoundError(err) {
			err = t.syncBatchWithLock(stateMutex, batch)
			if err != nil {
				return err
			}
		} else {
			err = t.syncExistingBatch(stateMutex, batch, localBatch)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *transactionExecutor) syncExistingBatch(mutex *sync.Mutex, batch *eth.DecodedBatch, localBatch *models.Batch) error {
	if batch.TransactionHash == localBatch.TransactionHash {
		batch.ID = localBatch.ID
		err := t.storage.MarkBatchAsSubmitted(&batch.Batch)
		if err != nil {
			return err
		}

		log.Printf(
			"Synced new existing batch. Batch number: %d. Batch Hash: %v",
			batch.Number.Uint64(),
			batch.Hash,
		)
	} else {
		txSender, err := t.getTransactionSender(batch.TransactionHash)
		if err != nil {
			return err
		}
		if *txSender != t.client.ChainConnection.GetAccount().From {
			return t.revertBatches(mutex, batch, localBatch)
		} else {
			// TODO our previous transaction must have failed this should never happen
			return ErrBatchSubmissionFailed
		}
	}
	return nil
}

func (t *transactionExecutor) revertBatches(mutex *sync.Mutex, batch *eth.DecodedBatch, localBatch *models.Batch) error {
	mutex.Lock()
	defer mutex.Unlock()

	stateTree := st.NewStateTree(t.storage)
	err := stateTree.RevertTo(*localBatch.PrevStateRoot)
	if err != nil {
		return err
	}
	err = t.revertBatchesInRange(&batch.Number)
	if err != nil {
		return err
	}
	return t.syncBatch(batch)
}

func (t *transactionExecutor) revertBatchesInRange(startNumber *models.Uint256) error {
	batches, err := t.storage.GetBatchesInRange(startNumber, nil)
	if err != nil {
		return err
	}
	batchIDs := make([]int32, 0, len(batches))
	for i := range batches {
		batchIDs = append(batchIDs, batches[i].ID)
	}
	err = t.excludeTransactionsFromCommitment(batchIDs...)
	if err != nil {
		return err
	}
	err = t.storage.DeleteCommitmentsByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	return t.storage.DeleteBatches(batchIDs...)
}

func (t *transactionExecutor) excludeTransactionsFromCommitment(batchIDs ...int32) error {
	hashes, err := t.storage.GetTransactionHashesByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, nil)
}

func (t *transactionExecutor) getTransactionSender(txHash common.Hash) (*common.Address, error) {
	tx, _, err := t.client.ChainConnection.GetBackend().TransactionByHash(t.ctx, txHash)
	if err != nil {
		return nil, err
	}
	message, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		return nil, err
	}
	sender := message.From()
	return &sender, nil
}

func getLatestBatchNumber(storage *st.Storage) (*models.Uint256, error) {
	latestBatch, err := storage.GetLatestSubmittedBatch()
	if st.IsNotFoundError(err) {
		return models.NewUint256(0), nil
	} else if err != nil {
		return nil, err
	}
	return &latestBatch.Number, nil
}

func (t *transactionExecutor) syncBatchWithLock(stateMutex *sync.Mutex, batch *eth.DecodedBatch) error {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	return t.syncBatch(batch)
}

func (t *transactionExecutor) syncBatch(batch *eth.DecodedBatch) error {
	batchID, err := t.storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

	batch.Batch.ID = *batchID

	switch batch.Type {
	case txtype.Transfer:
		err = t.syncTransferCommitments(batch)
		if err != nil {
			return err
		}
	case txtype.Create2Transfer:
		err = t.syncCreate2TransferCommitments(batch)
		if err != nil {
			return err
		}
	case txtype.MassMigration:
		return fmt.Errorf("unsupported batch type for sync: %s", batch.Type)
	}

	log.Printf("Synced new batch #%s from chain: %d commitments included", batch.Number.String(), len(batch.Commitments))
	return nil
}
