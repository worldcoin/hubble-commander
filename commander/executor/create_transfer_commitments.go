package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotEnoughTransfers = NewRollupError("not enough transfers")
)

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (c *RollupContext) CreateTransferCommitments(
	domain *bls.Domain,
) (commitments []models.Commitment, err error) {
	pendingTransfers, err := c.queryPendingTransfers()
	if err != nil {
		return nil, err
	}

	commitmentID, err := c.createCommitmentID()
	if err != nil {
		return nil, err
	}

	commitments = make([]models.Commitment, 0, c.cfg.MaxCommitmentsPerBatch)

	for i := uint8(0); len(commitments) != int(c.cfg.MaxCommitmentsPerBatch); i++ {
		var commitment *models.Commitment
		commitmentID.IndexInBatch = i

		pendingTransfers, commitment, err = c.createTransferCommitment(pendingTransfers.ToTransferArray(), commitmentID, domain)
		if err == ErrNotEnoughTransfers {
			break
		}
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
	}

	if len(commitments) == 0 {
		return nil, ErrNotEnoughTransfers
	}

	return commitments, nil
}

func (c *RollupContext) createTransferCommitment(
	pendingTransfers models.GenericTransactionArray,
	commitmentID *models.CommitmentID,
	domain *bls.Domain,
) (
	newPendingTransfers models.GenericTransactionArray,
	commitment *models.Commitment,
	err error,
) {
	startTime := time.Now()

	pendingTransfers, err = c.refillPendingTransfers(pendingTransfers)
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

	appliedTransfers, newPendingTransfers, err := c.applyTransfersForCommitment(pendingTransfers.ToTransferArray(), feeReceiver)
	if err == ErrNotEnoughTransfers {
		if revertErr := c.storage.StateTree.RevertTo(*initialStateRoot); revertErr != nil {
			return nil, nil, revertErr
		}
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	commitment, err = c.buildTransferCommitment(appliedTransfers.ToTransferArray(), commitmentID, feeReceiver.StateID, domain)
	if err != nil {
		return nil, nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		txtype.Transfer,
		appliedTransfers.Len(),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return newPendingTransfers, commitment, nil
}

func (c *RollupContext) applyTransfersForCommitment(pendingTransfers models.GenericTransactionArray, feeReceiver *FeeReceiver) (
	appliedTransfers, newPendingTransfers models.GenericTransactionArray,
	err error,
) {
	appliedTransfers = c.Executor.makeTransactionArray(0, c.cfg.MaxTxsPerCommitment)
	invalidTransfers := c.Executor.makeTransactionArray(0, 1)

	for {
		var transfers *AppliedTransfers

		numNeededTransfers := c.cfg.MaxTxsPerCommitment - uint32(appliedTransfers.Len())
		transfers, err = c.ApplyTransfers(pendingTransfers.ToTransferArray(), numNeededTransfers, feeReceiver)
		if err != nil {
			return nil, nil, err
		}

		appliedTransfers = appliedTransfers.Append(transfers.appliedTransfers)
		invalidTransfers = invalidTransfers.Append(transfers.invalidTransfers)

		if appliedTransfers.Len() == int(c.cfg.MaxTxsPerCommitment) {
			newPendingTransfers = removeTransfers(pendingTransfers, appliedTransfers.Append(invalidTransfers))
			return appliedTransfers, newPendingTransfers, nil
		}

		morePendingTransfers, err := c.queryMorePendingTransfers(appliedTransfers)
		if err == ErrNotEnoughTransfers {
			newPendingTransfers = removeTransfers(pendingTransfers, appliedTransfers.Append(invalidTransfers))
			return appliedTransfers, newPendingTransfers, nil
		}
		if err != nil {
			return nil, nil, err
		}
		pendingTransfers = morePendingTransfers
	}
}

func (c *RollupContext) refillPendingTransfers(pendingTransfers models.GenericTransactionArray) (models.GenericTransactionArray, error) {
	if pendingTransfers.Len() < int(c.cfg.MaxTxsPerCommitment) {
		return c.queryPendingTransfers()
	}
	return pendingTransfers, nil
}

func (c *RollupContext) queryPendingTransfers() (models.GenericTransactionArray, error) {
	pendingTransfers, err := c.Executor.getPendingTransactions(c.cfg.MaxCommitmentsPerBatch * c.cfg.MaxTxsPerCommitment)
	if err != nil {
		return nil, err
	}
	if pendingTransfers.Len() < int(c.cfg.MinTxsPerCommitment) {
		return nil, ErrNotEnoughTransfers
	}
	return pendingTransfers, nil
}

func (c *RollupContext) queryMorePendingTransfers(appliedTransfers models.GenericTransactionArray) (models.GenericTransactionArray, error) {
	numAppliedTransfers := uint32(appliedTransfers.Len())
	pendingTransfers, err := c.Executor.getPendingTransactions(
		c.cfg.MaxCommitmentsPerBatch*c.cfg.MaxTxsPerCommitment + numAppliedTransfers,
	)
	if err != nil {
		return nil, err
	}
	pendingTransfers = removeTransfers(pendingTransfers, appliedTransfers)

	if pendingTransfers.Len() < int(c.cfg.MinTxsPerCommitment) {
		return nil, ErrNotEnoughTransfers
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
