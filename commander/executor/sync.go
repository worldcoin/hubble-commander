package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrBatchSubmissionFailed = errors.New("previous submit batch transaction failed")
)

func (t *TransactionExecutor) SyncBatch(remoteBatch *eth.DecodedBatch) error {
	localBatch, err := t.storage.GetBatch(remoteBatch.ID)
	if err != nil && !st.IsNotFoundError(err) {
		return err
	}

	if st.IsNotFoundError(err) {
		return t.syncNewBatch(remoteBatch)
	} else {
		return t.syncExistingBatch(remoteBatch, localBatch)
	}
}

func (t *TransactionExecutor) syncExistingBatch(remoteBatch *eth.DecodedBatch, localBatch *models.Batch) error {
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
		return ErrBatchSubmissionFailed
	}
	return nil
}

func (t *TransactionExecutor) RevertBatches(startBatch *models.Batch) error {
	err := t.storage.StateTree.RevertTo(*startBatch.PrevStateRoot)
	if err != nil {
		return err
	}
	return t.revertBatchesFrom(&startBatch.ID)
}

func (t *TransactionExecutor) revertBatchesFrom(startBatchID *models.Uint256) error {
	batches, err := t.storage.GetBatchesInRange(startBatchID, nil)
	if err != nil {
		return err
	}
	numBatches := len(batches)
	batchIDs := make([]models.Uint256, 0, numBatches)
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
	log.Debugf("Removing %d local batches", numBatches)
	return t.storage.DeleteBatches(batchIDs...)
}

func (t *TransactionExecutor) excludeTransactionsFromCommitment(batchIDs ...models.Uint256) error {
	hashes, err := t.storage.GetTransactionHashesByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, nil)
}

func (t *TransactionExecutor) syncNewBatch(batch *eth.DecodedBatch) error {
	numCommitments := len(batch.Commitments)
	log.Debugf("Syncing new batch #%s with %d commitment(s) from chain", batch.ID.String(), numCommitments)
	err := t.storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

	err = t.syncCommitments(batch)
	if err != nil {
		return err
	}

	log.Printf("Synced new batch #%s with %d commitment(s) from chain", batch.ID.String(), numCommitments)
	return nil
}

func (t *TransactionExecutor) syncCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		log.WithFields(log.Fields{"batchID": batch.ID.String()}).Debugf("Syncing commitment #%d", i+1)
		err := t.syncCommitment(batch, &batch.Commitments[i])

		var disputableErr *DisputableError
		if errors.As(err, &disputableErr) {
			return disputableErr.WithCommitmentIndex(i)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TransactionExecutor) syncCommitment(
	batch *eth.DecodedBatch,
	commitment *encoder.DecodedCommitment,
) error {
	if len(commitment.Transactions)%encoder.GetTransactionLength(batch.Type) != 0 {
		return ErrInvalidDataLength
	}

	var transactions models.GenericTransactionArray
	var err error
	switch batch.Type {
	case txtype.Transfer:
		transactions, err = t.syncTransferCommitment(commitment)
	case txtype.Create2Transfer:
		transactions, err = t.syncCreate2TransferCommitment(commitment)
	case txtype.Genesis, txtype.MassMigration:
		return errors.Errorf("unsupported batch type for sync: %s", batch.Type)
	}
	if err != nil {
		return err
	}

	commitmentID, err := t.storage.AddCommitment(&models.Commitment{
		Type:              batch.Type,
		Transactions:      commitment.Transactions,
		FeeReceiver:       commitment.FeeReceiver,
		CombinedSignature: commitment.CombinedSignature,
		PostStateRoot:     commitment.StateRoot,
		IncludedInBatch:   &batch.ID,
	})
	if err != nil {
		return err
	}
	for i := 0; i < transactions.Len(); i++ {
		transactions.At(i).GetBase().IncludedInCommitment = commitmentID
	}

	for i := 0; i < transactions.Len(); i++ {
		hashTransfer, err := encoder.HashGenericTransaction(transactions.At(i))
		if err != nil {
			return err
		}
		transactions.At(i).GetBase().Hash = *hashTransfer
	}

	if transactions.Len() == 0 {
		return nil
	}
	return t.storage.BatchAddGenericTransaction(transactions)
}
