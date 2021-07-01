package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var (
	ErrFraudulentTransfer    = errors.New("fraudulent transfer encountered when syncing")
	ErrTransfersNotApplied   = errors.New("could not apply all transfers from synced batch") // TODO-AFS get rid of this error
	ErrBatchSubmissionFailed = errors.New("previous submit batch transaction failed")
	ErrInvalidSignature      = errors.New("invalid commitment signature")
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

func (t *TransactionExecutor) revertBatches(remoteBatch *eth.DecodedBatch, localBatch *models.Batch) error {
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

func (t *TransactionExecutor) revertBatchesInRange(startBatchID *models.Uint256) error {
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

func (t *TransactionExecutor) excludeTransactionsFromCommitment(batchIDs ...models.Uint256) error {
	hashes, err := t.storage.GetTransactionHashesByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, nil)
}

func (t *TransactionExecutor) getTransactionSender(txHash common.Hash) (*common.Address, error) {
	tx, _, err := t.client.ChainConnection.GetBackend().TransactionByHash(t.opts.Ctx, txHash)
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

func (t *TransactionExecutor) syncNewBatch(batch *eth.DecodedBatch) error {
	err := t.storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

	err = t.syncCommitments(batch)
	if err != nil {
		return err
	}

	log.Printf("Synced new batch #%s from chain with %d commitment(s)", batch.ID.String(), len(batch.Commitments))
	return nil
}

func (t *TransactionExecutor) syncCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		err := t.syncCommitment(batch, &batch.Commitments[i])
		if err == ErrInvalidSignature {
			// TODO: dispute fraudulent commitment
			return err
		}
		if IsDisputableTransferError(err) {
			// TODO: dispute fraudulent commitment
			return err
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
	case txtype.MassMigration:
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

	return t.storage.BatchAddGenericTransaction(transactions)
}
