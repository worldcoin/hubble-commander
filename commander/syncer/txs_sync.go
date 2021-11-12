package syncer

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *TxsContext) SyncCommitments(remoteBatch eth.DecodedBatch) error {
	batch := remoteBatch.ToDecodedTxBatch()
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

func (c *TxsContext) UpdateExistingBatch(batch eth.DecodedBatch) error {
	txBatch := batch.ToDecodedTxBatch()
	err := c.storage.UpdateBatch(&txBatch.Batch)
	if err != nil {
		return err
	}
	return c.setCommitmentsBodyHash(txBatch)
}

func (c *TxsContext) setCommitmentsBodyHash(batch *eth.DecodedTxBatch) error {
	commitments, err := c.storage.GetTxCommitmentsByBatchID(batch.ID)
	if err != nil {
		return err
	}
	for i := range commitments {
		commitments[i].BodyHash = batch.Commitments[i].BodyHash(*batch.AccountTreeRoot)
	}

	return c.storage.UpdateCommitments(commitments)
}

func (c *TxsContext) syncCommitment(batch *eth.DecodedTxBatch, commitment *encoder.DecodedCommitment) error {
	transactions, err := c.syncTxCommitment(commitment)
	if err != nil {
		return err
	}

	err = c.addCommitment(batch, commitment)
	if err != nil {
		return err
	}

	// TODO move adding transactions to syncTxCommitment
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

func (c *TxsContext) addCommitment(batch *eth.DecodedTxBatch, decodedCommitment *encoder.DecodedCommitment) error {
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
