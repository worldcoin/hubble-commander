package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrBatchSubmissionFailed = errors.New("previous submit batch transaction failed")
)

func (c *SyncContext) SyncBatch(remoteBatch *eth.DecodedBatch) error {
	localBatch, err := c.storage.GetBatch(remoteBatch.ID)
	if err != nil && !st.IsNotFoundError(err) {
		return err
	}

	if st.IsNotFoundError(err) {
		return c.syncNewBatch(remoteBatch)
	} else {
		return c.syncExistingBatch(remoteBatch, localBatch)
	}
}

func (c *SyncContext) syncExistingBatch(remoteBatch *eth.DecodedBatch, localBatch *models.Batch) error {
	if remoteBatch.TransactionHash == localBatch.TransactionHash {
		err := c.storage.MarkBatchAsSubmitted(&remoteBatch.Batch)
		if err != nil {
			return err
		}

		log.Printf(
			"Synced new existing batch. Batch ID: %d. Batch Hash: %v",
			remoteBatch.ID.Uint64(),
			remoteBatch.Hash,
		)
	} else {
		txSender, err := c.getTransactionSender(remoteBatch.TransactionHash)
		if err != nil {
			return err
		}
		if *txSender != c.client.ChainConnection.GetAccount().From {
			return NewInconsistentBatchError(localBatch)
		} else {
			// TODO remove the above check and this error once we use contracts with batchID verification:
			//  https://github.com/thehubbleproject/hubble-contracts/pull/601
			return ErrBatchSubmissionFailed
		}
	}
	return nil
}

func (c *SyncContext) getTransactionSender(txHash common.Hash) (*common.Address, error) {
	tx, _, err := c.client.ChainConnection.GetBackend().TransactionByHash(c.ctx, txHash)
	if err != nil {
		return nil, err
	}
	sender, err := types.LatestSignerForChainID(tx.ChainId()).Sender(tx)
	if err != nil {
		return nil, err
	}
	return &sender, nil
}

func (c *SyncContext) syncNewBatch(batch *eth.DecodedBatch) error {
	numCommitments := len(batch.Commitments)
	log.Debugf("Syncing new batch #%s with %d commitment(s) from chain", batch.ID.String(), numCommitments)
	err := c.storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

	err = c.syncCommitments(batch)
	if err != nil {
		return err
	}

	log.Printf("Synced new batch #%s with %d commitment(s) from chain", batch.ID.String(), numCommitments)
	return nil
}

func (c *SyncContext) syncCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		log.WithFields(log.Fields{"batchID": batch.ID.String()}).Debugf("Syncing commitment #%d", i+1)
		err := c.syncCommitment(batch, &batch.Commitments[i])

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

func (c *SyncContext) syncCommitment(
	batch *eth.DecodedBatch,
	commitment *encoder.DecodedCommitment,
) error {
	if len(commitment.Transactions)%encoder.GetTransactionLength(batch.Type) != 0 {
		return ErrInvalidDataLength
	}

	var transactions models.GenericTransactionArray
	var err error
	switch batch.Type {
	case batchtype.Transfer:
		transactions, err = c.syncTransferCommitment(commitment)
	case batchtype.Create2Transfer:
		transactions, err = c.syncCreate2TransferCommitment(commitment)
	case batchtype.Genesis, batchtype.MassMigration, batchtype.Deposit:
		return errors.Errorf("unsupported batch type for sync: %s", batch.Type)
	}
	if err != nil {
		return err
	}

	err = c.storage.AddCommitment(&models.Commitment{
		ID:                commitment.ID,
		Type:              batch.Type,
		Transactions:      commitment.Transactions,
		FeeReceiver:       commitment.FeeReceiver,
		CombinedSignature: commitment.CombinedSignature,
		PostStateRoot:     commitment.StateRoot,
	})
	if err != nil {
		return err
	}
	for i := 0; i < transactions.Len(); i++ {
		transactions.At(i).GetBase().CommitmentID = &commitment.ID
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
	return c.storage.BatchAddGenericTransaction(transactions)
}
