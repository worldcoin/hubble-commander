package executor

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

var pendingTxsCountMultiplier = uint32(2)

func (t *TransactionExecutor) CreateAndSubmitBatch(batchType txtype.TransactionType, domain *bls.Domain) (err error) {
	startTime := time.Now()
	var commitments []models.Commitment
	batch, err := t.NewPendingBatch(batchType)
	if err != nil {
		return err
	}

	if batchType == txtype.Transfer {
		commitments, err = t.buildTransferCommitments(domain)
	} else {
		commitments, err = t.buildCreate2TransfersCommitments(domain)
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

func (t *TransactionExecutor) buildTransferCommitments(domain *bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := t.storage.GetPendingTransfers(pendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
	if err != nil {
		return nil, err
	}
	return t.CreateTransferCommitments(pendingTransfers, domain)
}

func (t *TransactionExecutor) buildCreate2TransfersCommitments(domain *bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := t.storage.GetPendingCreate2Transfers(pendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
	if err != nil {
		return nil, err
	}
	return t.CreateCreate2TransferCommitments(pendingTransfers, domain)
}

func (t *TransactionExecutor) NewPendingBatch(batchType txtype.TransactionType) (*models.Batch, error) {
	stateTree := st.NewStateTree(t.storage)
	prevStateRoot, err := stateTree.Root()
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
