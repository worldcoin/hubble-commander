package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotEnoughTxs = NewRollupError("not enough transactions")
)

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (c *RollupContext) CreateCommitments() ([]models.Commitment, error) {
	pendingTxs, err := c.queryPendingTxs()
	if err != nil {
		return nil, err
	}

	commitmentID, err := c.nextCommitmentID()
	if err != nil {
		return nil, err
	}

	commitments := make([]models.Commitment, 0, c.cfg.MaxCommitmentsPerBatch)

	for i := uint8(0); len(commitments) != int(c.cfg.MaxCommitmentsPerBatch); i++ {
		var commitment *models.Commitment
		commitmentID.IndexInBatch = i

		pendingTxs, commitment, err = c.createCommitment(pendingTxs, commitmentID)
		if err == ErrNotEnoughTxs {
			break
		}
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
	}

	if len(commitments) == 0 {
		return nil, ErrNotEnoughTxs
	}

	return commitments, nil
}

func (c *RollupContext) createCommitment(pendingTxs models.GenericTransactionArray, commitmentID *models.CommitmentID) (
	newPendingTxs models.GenericTransactionArray,
	commitment *models.Commitment,
	err error,
) {
	startTime := time.Now()

	pendingTxs, err = c.refillPendingTxs(pendingTxs)
	if err != nil {
		return nil, nil, err
	}

	feeReceiver, err := c.getCommitmentFeeReceiver()
	if err != nil {
		return nil, nil, err
	}

	initialStateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, nil, err
	}

	applyResult, newPendingTxs, err := c.applyTxsForCommitment(pendingTxs, feeReceiver)
	if err == ErrNotEnoughTxs {
		if revertErr := c.storage.StateTree.RevertTo(*initialStateRoot); revertErr != nil {
			return nil, nil, revertErr
		}
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	commitment, err = c.buildCommitment(applyResult, commitmentID, feeReceiver.StateID)
	if err != nil {
		return nil, nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		c.BatchType,
		applyResult.AppliedTxs().Len(),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return newPendingTxs, commitment, nil
}

func (c *RollupContext) applyTxsForCommitment(pendingTxs models.GenericTransactionArray, feeReceiver *FeeReceiver) (
	result ApplyTxsForCommitmentResult,
	newPendingTxs models.GenericTransactionArray,
	err error,
) {
	aggregateResult := c.Executor.NewApplyTxsResult(c.cfg.MaxTxsPerCommitment)

	for {
		numNeededTxs := c.cfg.MaxTxsPerCommitment - uint32(aggregateResult.AppliedTxs().Len())
		applyTxsResult, err := c.ApplyTxs(pendingTxs, numNeededTxs, feeReceiver)
		if err != nil {
			return nil, nil, err
		}

		aggregateResult.AddApplyResult(applyTxsResult)

		if aggregateResult.AppliedTxs().Len() == int(c.cfg.MaxTxsPerCommitment) {
			newPendingTxs = removeTxs(pendingTxs, aggregateResult.AllTxs())
			return c.Executor.NewApplyTxsForCommitmentResult(aggregateResult), newPendingTxs, nil
		}

		morePendingTransfers, err := c.queryMorePendingTxs(aggregateResult.AppliedTxs())
		if err == ErrNotEnoughTxs {
			newPendingTxs = removeTxs(pendingTxs, aggregateResult.AllTxs())
			return c.Executor.NewApplyTxsForCommitmentResult(aggregateResult), newPendingTxs, nil
		}
		if err != nil {
			return nil, nil, err
		}
		pendingTxs = morePendingTransfers
	}
}

func (c *RollupContext) refillPendingTxs(pendingTxs models.GenericTransactionArray) (models.GenericTransactionArray, error) {
	if pendingTxs.Len() < int(c.cfg.MaxTxsPerCommitment) {
		return c.queryPendingTxs()
	}
	return pendingTxs, nil
}

func (c *RollupContext) queryPendingTxs() (models.GenericTransactionArray, error) {
	pendingTxs, err := c.Executor.GetPendingTxs(c.cfg.MaxCommitmentsPerBatch * c.cfg.MaxTxsPerCommitment)
	if err != nil {
		return nil, err
	}
	if pendingTxs.Len() < int(c.cfg.MinTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}
	return pendingTxs, nil
}

func (c *RollupContext) queryMorePendingTxs(appliedTxs models.GenericTransactionArray) (models.GenericTransactionArray, error) {
	numAppliedTransfers := uint32(appliedTxs.Len())
	pendingTransfers, err := c.Executor.GetPendingTxs(
		c.cfg.MaxCommitmentsPerBatch*c.cfg.MaxTxsPerCommitment + numAppliedTransfers,
	)
	if err != nil {
		return nil, err
	}
	pendingTransfers = removeTxs(pendingTransfers, appliedTxs)

	if pendingTransfers.Len() < int(c.cfg.MinTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}
	return pendingTransfers, nil
}

func (c *RollupContext) getCommitmentFeeReceiver() (*FeeReceiver, error) {
	commitmentTokenID := models.MakeUint256(0) // TODO support multiple tokens
	feeReceiverState, err := c.storage.GetFeeReceiverStateLeaf(c.cfg.FeeReceiverPubKeyID, commitmentTokenID)
	if err != nil {
		return nil, err
	}
	return &FeeReceiver{
		StateID: feeReceiverState.StateID,
		TokenID: feeReceiverState.TokenID,
	}, nil
}

func removeTxs(txList, toRemove models.GenericTransactionArray) models.GenericTransactionArray {
	outputIndex := 0
	for i := 0; i < txList.Len(); i++ {
		tx := txList.At(i)
		if !txExists(toRemove, tx) {
			txList.Set(outputIndex, tx)
			outputIndex++
		}
	}

	return txList.Slice(0, outputIndex)
}

func txExists(txList models.GenericTransactionArray, tx models.GenericTransaction) bool {
	for i := 0; i < txList.Len(); i++ {
		if txList.At(i).GetBase().Hash == tx.GetBase().Hash {
			return true
		}
	}
	return false
}