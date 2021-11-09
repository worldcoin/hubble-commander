package syncer

import (
	"context"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrBatchSubmissionFailed = errors.New("previous submit batch transaction failed")
)

func (c *Context) SyncBatch(batch eth.DecodedBatch) error {
	remoteBatch := batch.ToDecodedTxBatch()
	localBatch, err := c.storage.GetBatch(remoteBatch.ID)
	if err != nil && !st.IsNotFoundError(err) {
		return err
	}

	if st.IsNotFoundError(err) {
		return c.SyncNewBatch(remoteBatch)
	} else {
		return c.syncExistingBatch(remoteBatch, localBatch)
	}
}

func (c *Context) syncExistingBatch(remoteBatch *eth.DecodedTxBatch, localBatch *models.Batch) error {
	if remoteBatch.TransactionHash == localBatch.TransactionHash {
		err := c.UpdateExistingBatch(remoteBatch)
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
		if *txSender != c.client.Blockchain.GetAccount().From {
			return NewInconsistentBatchError(localBatch)
		} else {
			// TODO remove the above check and this error once we use contracts with batchID verification:
			//  https://github.com/thehubbleproject/hubble-contracts/pull/601
			return ErrBatchSubmissionFailed
		}
	}
	return nil
}

func (c *Context) UpdateExistingBatch(batch eth.DecodedBatch) error {
	txBatch := batch.ToDecodedTxBatch()
	err := c.storage.UpdateBatch(&txBatch.Batch)
	if err != nil {
		return err
	}
	return c.setCommitmentsBodyHash(txBatch)
}

func (c *Context) getTransactionSender(txHash common.Hash) (*common.Address, error) {
	tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, err
	}
	sender, err := types.LatestSignerForChainID(tx.ChainId()).Sender(tx)
	if err != nil {
		return nil, err
	}
	return &sender, nil
}

func (c *Context) setCommitmentsBodyHash(batch *eth.DecodedTxBatch) error {
	commitments, err := c.storage.GetTxCommitmentsByBatchID(batch.ID)
	if err != nil {
		return err
	}
	for i := range commitments {
		commitments[i].BodyHash = batch.Commitments[i].BodyHash(*batch.AccountTreeRoot)
	}

	return c.storage.UpdateCommitments(commitments)
}

func (c *Context) SyncNewBatch(remoteBatch eth.DecodedBatch) error {
	batch := remoteBatch.ToDecodedTxBatch()
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

func (c *Context) syncCommitments(batch *eth.DecodedTxBatch) error {
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

func (c *Context) syncCommitment(batch *eth.DecodedTxBatch, commitment *encoder.DecodedCommitment) error {
	transactions, err := c.syncTxCommitment(commitment)
	if err != nil {
		return err
	}

	err = c.addTxCommitment(batch, commitment)
	if err != nil {
		return err
	}
	for i := 0; i < transactions.Len(); i++ {
		transactions.At(i).GetBase().CommitmentID = &commitment.ID
		hashTransfer, err := c.Syncer.HashTx(transactions.At(i))
		if err != nil {
			return err
		}
		transactions.At(i).GetBase().Hash = *hashTransfer
	}

	if transactions.Len() == 0 {
		return nil
	}
	return c.Syncer.BatchAddTxs(transactions)
}

func (c *Context) addTxCommitment(batch *eth.DecodedTxBatch, decodedCommitment *encoder.DecodedCommitment) error {
	commitment := &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID:            decodedCommitment.ID,
			Type:          batch.Type,
			PostStateRoot: decodedCommitment.StateRoot,
		},
		FeeReceiver:       decodedCommitment.FeeReceiver,
		CombinedSignature: decodedCommitment.CombinedSignature,
		BodyHash:          decodedCommitment.BodyHash(*batch.AccountTreeRoot),
	}

	return c.storage.AddTxCommitment(commitment)
}
