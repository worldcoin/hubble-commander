package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotEnoughC2Transfers = NewRollupError("not enough create2transfers")
)

func (c *RollupContext) CreateCreate2TransferCommitments(domain *bls.Domain) (commitments []models.Commitment, err error) {
	pendingTransfers, err := c.queryPendingC2Ts()
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

		pendingTransfers, commitment, err = c.createC2TCommitment(pendingTransfers, commitmentID, domain)
		if err == ErrNotEnoughC2Transfers {
			break
		}
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
	}

	if len(commitments) == 0 {
		return nil, ErrNotEnoughC2Transfers
	}

	return commitments, nil
}

func (c *RollupContext) createC2TCommitment(
	pendingTransfers []models.Create2Transfer,
	commitmentID *models.CommitmentID,
	domain *bls.Domain,
) (newPendingTransfers []models.Create2Transfer, commitment *models.Commitment, err error) {
	startTime := time.Now()

	pendingTransfers, err = c.refillPendingC2Ts(pendingTransfers)
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

	appliedTransfers, newPendingTransfers, addedPubKeyIDs, err := c.applyC2TsForCommitment(pendingTransfers, feeReceiver)
	if err == ErrNotEnoughC2Transfers {
		if revertErr := c.storage.StateTree.RevertTo(*initialStateRoot); revertErr != nil {
			return nil, nil, revertErr
		}
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	commitment, err = c.buildC2TCommitment(appliedTransfers, addedPubKeyIDs, commitmentID, feeReceiver.StateID, domain)
	if err != nil {
		return nil, nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		txtype.Create2Transfer,
		len(appliedTransfers),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return newPendingTransfers, commitment, nil
}

func (c *RollupContext) applyC2TsForCommitment(pendingTransfers []models.Create2Transfer, feeReceiver *FeeReceiver) (
	appliedTransfers, newPendingTransfers []models.Create2Transfer,
	addedPubKeyIDs []uint32,
	err error,
) {
	appliedTransfers = make([]models.Create2Transfer, 0, c.cfg.MaxTxsPerCommitment)
	invalidTransfers := make([]models.Create2Transfer, 0, 1)
	addedPubKeyIDs = make([]uint32, 0, c.cfg.MaxTxsPerCommitment)

	for {
		var transfers *AppliedC2Transfers

		numNeededTransfers := c.cfg.MaxTxsPerCommitment - uint32(len(appliedTransfers))
		transfers, err = c.ApplyCreate2Transfers(pendingTransfers, numNeededTransfers, feeReceiver)
		if err != nil {
			return nil, nil, nil, err
		}

		appliedTransfers = append(appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)
		addedPubKeyIDs = append(addedPubKeyIDs, transfers.addedPubKeyIDs...)

		if len(appliedTransfers) == int(c.cfg.MaxTxsPerCommitment) {
			newPendingTransfers = removeC2Ts(pendingTransfers, append(appliedTransfers, invalidTransfers...))
			return appliedTransfers, newPendingTransfers, addedPubKeyIDs, nil
		}

		morePendingTransfers, err := c.queryMorePendingC2Ts(appliedTransfers)
		if err == ErrNotEnoughC2Transfers {
			newPendingTransfers = removeC2Ts(pendingTransfers, append(appliedTransfers, invalidTransfers...))
			return appliedTransfers, newPendingTransfers, addedPubKeyIDs, nil
		}
		if err != nil {
			return nil, nil, nil, err
		}
		pendingTransfers = morePendingTransfers
	}
}

func (c *ExecutionContext) refillPendingC2Ts(pendingTransfers []models.Create2Transfer) ([]models.Create2Transfer, error) {
	if len(pendingTransfers) < int(c.cfg.MaxTxsPerCommitment) {
		return c.queryPendingC2Ts()
	}
	return pendingTransfers, nil
}

func (c *ExecutionContext) queryPendingC2Ts() ([]models.Create2Transfer, error) {
	pendingTransfers, err := c.storage.GetPendingCreate2Transfers(c.cfg.MaxCommitmentsPerBatch * c.cfg.MaxTxsPerCommitment)
	if err != nil {
		return nil, err
	}
	if len(pendingTransfers) < int(c.cfg.MinTxsPerCommitment) {
		return nil, ErrNotEnoughC2Transfers
	}
	return pendingTransfers, nil
}

func (c *ExecutionContext) queryMorePendingC2Ts(appliedTransfers []models.Create2Transfer) ([]models.Create2Transfer, error) {
	numAppliedTransfers := uint32(len(appliedTransfers))
	// TODO use SQL Offset instead
	pendingTransfers, err := c.storage.GetPendingCreate2Transfers(
		c.cfg.MaxCommitmentsPerBatch*c.cfg.MaxTxsPerCommitment + numAppliedTransfers,
	)
	if err != nil {
		return nil, err
	}
	pendingTransfers = removeC2Ts(pendingTransfers, appliedTransfers)

	numNeededTransfers := c.cfg.MaxTxsPerCommitment - numAppliedTransfers
	if len(pendingTransfers) < int(numNeededTransfers) {
		return nil, ErrNotEnoughC2Transfers
	}
	return pendingTransfers, nil
}

func removeC2Ts(transferList, toRemove []models.Create2Transfer) []models.Create2Transfer {
	outputIndex := 0
	for i := range transferList {
		transfer := &transferList[i]
		if !create2TransferExists(toRemove, transfer) {
			transferList[outputIndex] = *transfer
			outputIndex++
		}
	}

	return transferList[:outputIndex]
}

func create2TransferExists(transferList []models.Create2Transfer, tx *models.Create2Transfer) bool {
	for i := range transferList {
		if transferList[i].Hash == tx.Hash {
			return true
		}
	}
	return false
}
