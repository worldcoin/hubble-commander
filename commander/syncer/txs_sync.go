package syncer

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *TxsContext) SyncCommitments(remoteBatch eth.DecodedBatch) error {
	batch := remoteBatch.ToDecodedTxBatch()
	for i := range batch.Commitments {
		log.WithFields(log.Fields{"batchID": batch.ID.String()}).Debugf("Syncing commitment #%d", i+1)

		err := c.syncCommitment(batch, batch.Commitments[i])

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

func (c *TxsContext) UpdateExistingBatch(batch eth.DecodedBatch, prevStateRoot common.Hash) error {
	err := c.storage.UpdateBatch(batch.ToBatch(prevStateRoot))
	if err != nil {
		return err
	}
	return c.setCommitmentsBodyHash(batch.ToDecodedTxBatch())
}

func (c *TxsContext) syncCommitment(batch *eth.DecodedTxBatch, commitment encoder.Commitment) error {
	err := c.syncTxCommitment(commitment)
	if err != nil {
		return err
	}

	return c.addCommitment(batch, commitment)
}

func (c *TxsContext) addCommitment(batch *eth.DecodedTxBatch, encodedCommitment encoder.Commitment) (err error) {
	decodedCommitment := encodedCommitment.ToDecodedCommitment()

	var commitment models.Commitment
	switch batch.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		commitment = &models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				ID:            decodedCommitment.ID,
				Type:          batch.Type,
				PostStateRoot: decodedCommitment.StateRoot,
			},
			FeeReceiver:       decodedCommitment.FeeReceiver,
			CombinedSignature: decodedCommitment.CombinedSignature,
			BodyHash:          encodedCommitment.BodyHash(batch.AccountTreeRoot),
		}
	case batchtype.MassMigration:
		commitment = &models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				ID:            decodedCommitment.ID,
				Type:          batch.Type,
				PostStateRoot: decodedCommitment.StateRoot,
			},
			FeeReceiver:       decodedCommitment.FeeReceiver,
			CombinedSignature: decodedCommitment.CombinedSignature,
			BodyHash:          encodedCommitment.BodyHash(batch.AccountTreeRoot),
			Meta:              encodedCommitment.(*encoder.DecodedMMCommitment).Meta,
			WithdrawRoot:      encodedCommitment.(*encoder.DecodedMMCommitment).WithdrawRoot,
		}
	default:
		panic("invalid batch type")
	}

	return c.storage.AddCommitment(commitment)
}

func (c *TxsContext) setCommitmentsBodyHash(batch *eth.DecodedTxBatch) error {
	commitments, err := c.storage.GetCommitmentsByBatchID(batch.ID)
	if err != nil {
		return err
	}
	for i := range commitments {
		commitments[i].SetBodyHash(batch.Commitments[i].BodyHash(batch.AccountTreeRoot))
	}

	return c.storage.UpdateCommitments(commitments)
}
