package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	log "github.com/sirupsen/logrus"
)

func (t *ExecutionContext) CreateAndSubmitBatch(batchType txtype.TransactionType, domain *bls.Domain) (err error) {
	startTime := time.Now()
	var commitments []models.Commitment
	batch, err := t.NewPendingBatch(batchType)
	if err != nil {
		return err
	}

	if batchType == txtype.Transfer {
		commitments, err = t.CreateTransferCommitments(domain)
	} else {
		commitments, err = t.CreateCreate2TransferCommitments(domain)
	}
	if err != nil {
		return err
	}

	err = t.SubmitBatch(batch, commitments)
	if err != nil {
		return err
	}

	log.Printf(
		"Submitted a %s batch with %d commitment(s) on chain in %s. Batch ID: %d. Transaction hash: %v",
		batchType.String(),
		len(commitments),
		time.Since(startTime).Round(time.Millisecond).String(),
		batch.ID.Uint64(),
		batch.TransactionHash,
	)
	return nil
}

func (t *ExecutionContext) NewPendingBatch(batchType txtype.TransactionType) (*models.Batch, error) {
	prevStateRoot, err := t.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}
	batchID, err := t.storage.GetNextBatchID()
	if err != nil {
		return nil, err
	}
	return &models.Batch{
		ID:            *batchID,
		Type:          batchType,
		PrevStateRoot: prevStateRoot,
	}, nil
}
