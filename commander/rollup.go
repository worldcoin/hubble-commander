package commander

import (
	"errors"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (c *Commander) rollupLoop() (err error) {
	ticker := time.NewTicker(c.cfg.Rollup.BatchLoopInterval)
	defer ticker.Stop()

	currentBatchType := txtype.Transfer

	for {
		select {
		case <-c.stopChannel:
			return nil
		case <-ticker.C:
			err := c.rollupLoopIteration(&currentBatchType)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Commander) rollupLoopIteration(currentBatchType *txtype.TransactionType) (err error) {
	if !c.storage.IsProposer() {
		return nil
	}

	transactionExecutor, err := newTransactionExecutor(c.storage, c.client, c.cfg.Rollup)
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	if *currentBatchType == txtype.Transfer {
		err = transactionExecutor.CreateAndSubmitBatch(*currentBatchType, c.domain)
		*currentBatchType = txtype.Create2Transfer
	} else {
		err = transactionExecutor.CreateAndSubmitBatch(*currentBatchType, c.domain)
		*currentBatchType = txtype.Transfer
	}
	if err != nil {
		var e *RollupError
		if errors.As(err, &e) {
			return nil
		}
		return err
	}

	return transactionExecutor.Commit()
}

func (t *transactionExecutor) CreateAndSubmitBatch(batchType txtype.TransactionType, domain *bls.Domain) (err error) {
	var commitments []models.Commitment

	if batchType == txtype.Transfer {
		commitments, err = t.buildTransferCommitments(domain)
	} else {
		commitments, err = t.buildCreate2TransfersCommitments(domain)
	}
	if err != nil {
		return err
	}

	err = t.submitBatch(batchType, commitments)
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionExecutor) buildTransferCommitments(domain *bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := t.storage.GetPendingTransfers()
	if err != nil {
		return nil, err
	}
	return t.createTransferCommitments(pendingTransfers, domain)
}

func (t *transactionExecutor) buildCreate2TransfersCommitments(
	domain *bls.Domain,
) ([]models.Commitment, error) {
	pendingTransfers, err := t.storage.GetPendingCreate2Transfers()
	if err != nil {
		return nil, err
	}
	return t.createCreate2TransferCommitments(pendingTransfers, domain)
}
