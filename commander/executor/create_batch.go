package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	log "github.com/sirupsen/logrus"
)

func (c *RollupContext) CreateAndSubmitBatch(domain *bls.Domain) (err error) {
	startTime := time.Now()
	var commitments []models.Commitment
	batch, err := c.NewPendingBatch(c.BatchType)
	if err != nil {
		return err
	}

	if c.BatchType == txtype.Transfer {
		commitments, err = c.CreateTransferCommitments(domain)
	} else {
		commitments, err = c.CreateCreate2TransferCommitments(domain)
	}
	if err != nil {
		return err
	}

	err = c.SubmitBatch(batch, commitments)
	if err != nil {
		return err
	}

	log.Printf(
		"Submitted a %s batch with %d commitment(s) on chain in %s. Batch ID: %d. Transaction hash: %v",
		c.BatchType.String(),
		len(commitments),
		time.Since(startTime).Round(time.Millisecond).String(),
		batch.ID.Uint64(),
		batch.TransactionHash,
	)
	return nil
}

func (c *ExecutionContext) NewPendingBatch(batchType txtype.TransactionType) (*models.Batch, error) {
	prevStateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}
	batchID, err := c.storage.GetNextBatchID()
	if err != nil {
		return nil, err
	}
	return &models.Batch{
		ID:            *batchID,
		Type:          batchType,
		PrevStateRoot: prevStateRoot,
	}, nil
}
