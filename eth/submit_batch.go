package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SubmitBatchFunc func() (*types.Transaction, error)

func (c *Client) SubmitTxBatch(
	batchType batchtype.BatchType,
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (*types.Transaction, error) {
	// nolint:exhaustive
	switch batchType {
	case batchtype.Transfer:
		return c.SubmitTransfersBatch(batchID, commitments)
	case batchtype.Create2Transfer:
		return c.SubmitCreate2TransfersBatch(batchID, commitments)
	case batchtype.MassMigration:
		panic("not implemented")
	default:
		panic("invalid batch type")
	}
}

func (c *Client) SubmitTransfersBatch(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*types.Transaction, error) {
	return c.rollup().
		WithValue(c.config.StakeAmount).
		WithGasLimit(*c.config.TransferBatchSubmissionGasLimit).
		SubmitTransfer(encoder.CommitmentsToSubmitBatchFields(batchID, commitments))
}

func (c *Client) SubmitCreate2TransfersBatch(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*types.Transaction, error) {
	return c.rollup().
		WithValue(c.config.StakeAmount).
		WithGasLimit(*c.config.C2TBatchSubmissionGasLimit).
		SubmitCreate2Transfer(encoder.CommitmentsToSubmitBatchFields(batchID, commitments))
}

func (c *Client) SubmitTransfersBatchAndWait(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*models.Batch, error) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitTransfersBatch(batchID, commitments)
	})
}
func (c *Client) SubmitCreate2TransfersBatchAndWait(
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (*models.Batch, error) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitCreate2TransfersBatch(batchID, commitments)
	})
}

func (c *Client) submitBatchAndWait(submit SubmitBatchFunc) (batch *models.Batch, err error) {
	tx, err := submit()
	if err != nil {
		return
	}

	receipt, err := chain.WaitToBeMined(c.Blockchain.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log, err := retrieveLog(receipt, NewBatchEvent)
	if err != nil {
		return nil, err
	}

	event := new(rollup.RollupNewBatch)
	err = c.rollupContract.UnpackLog(event, NewBatchEvent, *log)
	if err != nil {
		return nil, err
	}

	return c.handleNewBatchEvent(event)
}

func (c *Client) handleNewBatchEvent(event *rollup.RollupNewBatch) (*models.Batch, error) {
	batch, err := c.GetBatch(models.NewUint256FromBig(*event.BatchID))
	if err != nil {
		return nil, err
	}
	batch.AccountTreeRoot = ref.Hash(common.BytesToHash(event.AccountRoot[:]))
	return batch, nil
}
