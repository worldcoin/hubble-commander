package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type SubmitBatchFunc func() (*types.Transaction, error)

func (c *Client) SubmitTransfersBatch(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*types.Transaction, error) {
	tx, err := c.rollup().
		WithValue(c.config.StakeAmount).
		WithGasLimit(*c.config.TransferBatchSubmissionGasLimit).
		SubmitTransfer(encoder.CommitmentsToTransferAndC2TSubmitBatchFields(batchID, commitments))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.txsHashesChan <- tx.Hash()

	return tx, nil
}

func (c *Client) SubmitCreate2TransfersBatch(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*types.Transaction, error) {
	tx, err := c.rollup().
		WithValue(c.config.StakeAmount).
		WithGasLimit(*c.config.C2TBatchSubmissionGasLimit).
		SubmitCreate2Transfer(encoder.CommitmentsToTransferAndC2TSubmitBatchFields(batchID, commitments))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.txsHashesChan <- tx.Hash()

	return tx, nil
}

func (c *Client) SubmitMassMigrationsBatch(
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (*types.Transaction, error) {
	tx, err := c.rollup().
		WithValue(c.config.StakeAmount).
		WithGasLimit(*c.config.MMBatchSubmissionGasLimit).
		SubmitMassMigration(encoder.CommitmentsToSubmitMMBatchFields(batchID, commitments))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.txsHashesChan <- tx.Hash()

	return tx, nil
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

func (c *Client) SubmitMassMigrationsBatchAndWait(
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (*models.Batch, error) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitMassMigrationsBatch(batchID, commitments)
	})
}

func (c *Client) submitBatchAndWait(submit SubmitBatchFunc) (batch *models.Batch, err error) {
	tx, err := submit()
	if err != nil {
		return
	}

	receipt, err := c.WaitToBeMined(tx)
	if err != nil {
		return nil, err
	}

	log, err := retrieveLog(receipt, NewBatchEvent)
	if err != nil {
		return nil, err
	}

	event := new(rollup.RollupNewBatch)
	err = c.Rollup.BoundContract.UnpackLog(event, NewBatchEvent, *log)
	if err != nil {
		return nil, err
	}

	return c.handleNewBatchEvent(event)
}

func (c *Client) handleNewBatchEvent(event *rollup.RollupNewBatch) (*models.Batch, error) {
	contractBatch, err := c.GetContractBatch(models.NewUint256FromBig(*event.BatchID))
	if err != nil {
		return nil, err
	}
	batch := contractBatch.ToModelBatch()
	batch.AccountTreeRoot = ref.Hash(common.BytesToHash(event.AccountRoot[:]))
	return batch, nil
}
