package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotEnoughTxs = NewRollupError("not enough transactions")
)

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (c *RollupContext) CreateTxCommitments(domain *bls.Domain) ([]models.Commitment, error) {
	pendingTxs, err := c.queryPendingTxs()
	if err != nil {
		return nil, err
	}

	commitmentID, err := c.createCommitmentID()
	if err != nil {
		return nil, err
	}

	commitments := make([]models.Commitment, 0, c.cfg.MaxCommitmentsPerBatch)

	for i := uint8(0); len(commitments) != int(c.cfg.MaxCommitmentsPerBatch); i++ {
		var commitment *models.Commitment
		commitmentID.IndexInBatch = i

		pendingTxs, commitment, err = c.createTxCommitment(pendingTxs, commitmentID, domain)
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

func (c *RollupContext) createTxCommitment(
	pendingTxs models.GenericTransactionArray,
	commitmentID *models.CommitmentID,
	domain *bls.Domain,
) (
	newPendingTxs models.GenericTransactionArray,
	commitment *models.Commitment,
	err error,
) {
	startTime := time.Now()

	pendingTxs, err = c.refillPendingTransfers(pendingTxs)
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

	applyResult, newPendingTxs, err := c.applyTransfersForCommitment(pendingTxs, feeReceiver)
	if err == ErrNotEnoughTxs {
		if revertErr := c.storage.StateTree.RevertTo(*initialStateRoot); revertErr != nil {
			return nil, nil, revertErr
		}
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	commitment, err = c.buildTransferCommitment(applyResult, commitmentID, feeReceiver.StateID, domain)
	if err != nil {
		return nil, nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		txtype.Transfer,
		applyResult.AppliedTransfers().Len(),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return newPendingTxs, commitment, nil
}

func (c *RollupContext) applyTransfersForCommitment(pendingTransfers models.GenericTransactionArray, feeReceiver *FeeReceiver) (
	result ApplyTxsForCommitmentResult,
	newPendingTransfers models.GenericTransactionArray,
	err error,
) {
	aggregateResult := c.Executor.NewApplyTxsResult(c.cfg.MaxTxsPerCommitment)

	for {
		numNeededTransfers := c.cfg.MaxTxsPerCommitment - uint32(aggregateResult.AppliedTxs().Len())
		applyTxsResult, err := c.ApplyTxs(pendingTransfers, numNeededTransfers, feeReceiver)
		if err != nil {
			return nil, nil, err
		}

		aggregateResult.AddTxs(applyTxsResult)

		if aggregateResult.AppliedTxs().Len() == int(c.cfg.MaxTxsPerCommitment) {
			newPendingTransfers = removeTransfers(pendingTransfers, aggregateResult.AllTxs())
			return c.Executor.NewApplyTxsForCommitmentResult(aggregateResult), newPendingTransfers, nil
		}

		morePendingTransfers, err := c.queryMorePendingTransfers(aggregateResult.AppliedTxs())
		if err == ErrNotEnoughTxs {
			newPendingTransfers = removeTransfers(pendingTransfers, aggregateResult.AllTxs())
			return c.Executor.NewApplyTxsForCommitmentResult(aggregateResult), newPendingTransfers, nil
		}
		if err != nil {
			return nil, nil, err
		}
		pendingTransfers = morePendingTransfers
	}
}

func (c *RollupContext) refillPendingTransfers(pendingTransfers models.GenericTransactionArray) (models.GenericTransactionArray, error) {
	if pendingTransfers.Len() < int(c.cfg.MaxTxsPerCommitment) {
		return c.queryPendingTxs()
	}
	return pendingTransfers, nil
}

func (c *RollupContext) queryPendingTxs() (models.GenericTransactionArray, error) {
	pendingTransfers, err := c.Executor.GetPendingTxs(c.cfg.MaxCommitmentsPerBatch * c.cfg.MaxTxsPerCommitment)
	if err != nil {
		return nil, err
	}
	if pendingTransfers.Len() < int(c.cfg.MinTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}
	return pendingTransfers, nil
}

func (c *RollupContext) queryMorePendingTransfers(appliedTransfers models.GenericTransactionArray) (models.GenericTransactionArray, error) {
	numAppliedTransfers := uint32(appliedTransfers.Len())
	pendingTransfers, err := c.Executor.GetPendingTxs(
		c.cfg.MaxCommitmentsPerBatch*c.cfg.MaxTxsPerCommitment + numAppliedTransfers,
	)
	if err != nil {
		return nil, err
	}
	pendingTransfers = removeTransfers(pendingTransfers, appliedTransfers)

	if pendingTransfers.Len() < int(c.cfg.MinTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}
	return pendingTransfers, nil
}

func (c *ExecutionContext) getCommitmentFeeReceiver() (*FeeReceiver, error) {
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

func removeTransfers(transferList, toRemove models.GenericTransactionArray) models.GenericTransactionArray {
	outputIndex := 0
	for i := 0; i < transferList.Len(); i++ {
		tx := transferList.At(i)
		if !transferExists(toRemove, tx) {
			transferList.Set(outputIndex, tx)
			outputIndex++
		}
	}

	return transferList.Slice(0, outputIndex)
}

func transferExists(transferList models.GenericTransactionArray, tx models.GenericTransaction) bool {
	for i := 0; i < transferList.Len(); i++ {
		if transferList.At(i).GetBase().Hash == tx.GetBase().Hash {
			return true
		}
	}
	return false
}
