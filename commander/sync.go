package commander

import (
	"fmt"
	"log"

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

func (t *transactionExecutor) SyncBatches(startBlock, endBlock uint64) error {
	latestBatchID, err := getLatestBatchID(t.storage)
	if err != nil {
		return err
	}

	newRemoteBatches, err := t.client.GetBatches(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	})
	if err != nil {
		return err
	}

	for i := range newRemoteBatches {
		remoteBatch := &newRemoteBatches[i]
		if remoteBatch.ID.Cmp(latestBatchID) <= 0 {
			continue
		}

		localBatch, err := t.storage.GetBatch(remoteBatch.ID)
		if err != nil && !st.IsNotFoundError(err) {
			return err
		}

		if st.IsNotFoundError(err) {
			err = t.syncNewBatch(remoteBatch)
			if err != nil {
				return err
			}
		} else {
			err = t.syncExistingBatch(remoteBatch, localBatch)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *transactionExecutor) syncExistingBatch(remoteBatch *eth.DecodedBatch, localBatch *models.Batch) error {
	if remoteBatch.TransactionHash == localBatch.TransactionHash {
		err := t.storage.MarkBatchAsSubmitted(&remoteBatch.Batch)
		if err != nil {
			return err
		}

		log.Printf(
			"Synced new existing batch. Batch ID: %d. Batch Hash: %v",
			remoteBatch.ID.Uint64(),
			remoteBatch.Hash,
		)
	} else {
		txSender, err := t.getTransactionSender(remoteBatch.TransactionHash)
		if err != nil {
			return err
		}
		if *txSender != t.client.ChainConnection.GetAccount().From {
			return t.revertBatches(remoteBatch, localBatch)
		} else {
			// TODO remove the above check and this error once we use contracts with batchID verification:
			//  https://github.com/thehubbleproject/hubble-contracts/pull/601
			return ErrBatchSubmissionFailed
		}
	}
	return nil
}

func (t *transactionExecutor) revertBatches(remoteBatch *eth.DecodedBatch, localBatch *models.Batch) error {
	stateTree := st.NewStateTree(t.storage)
	err := stateTree.RevertTo(*localBatch.PrevStateRoot)
	if err != nil {
		return err
	}
	err = t.revertBatchesInRange(&remoteBatch.ID)
	if err != nil {
		return err
	}
	return t.syncNewBatch(remoteBatch)
}

func (t *transactionExecutor) revertBatchesInRange(startBatchID *models.Uint256) error {
	batches, err := t.storage.GetBatchesInRange(startBatchID, nil)
	if err != nil {
		return err
	}
	batchIDs := make([]models.Uint256, 0, len(batches))
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

func (t *transactionExecutor) excludeTransactionsFromCommitment(batchIDs ...models.Uint256) error {
	hashes, err := t.storage.GetTransactionHashesByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, nil)
}

func (t *transactionExecutor) getTransactionSender(txHash common.Hash) (*common.Address, error) {
	tx, _, err := t.client.ChainConnection.GetBackend().TransactionByHash(t.opts.ctx, txHash)
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

func getLatestBatchID(storage *st.Storage) (*models.Uint256, error) {
	latestBatch, err := storage.GetLatestSubmittedBatch()
	if st.IsNotFoundError(err) {
		return models.NewUint256(0), nil
	} else if err != nil {
		return nil, err
	}
	return &latestBatch.ID, nil
}

func (t *transactionExecutor) syncNewBatch(batch *eth.DecodedBatch) error {
	err := t.storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

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

	log.Printf("Synced new batch #%s from chain: %d commitments included", batch.ID.String(), len(batch.Commitments))
	return nil
}
