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
	opts := c.transactOpts(c.config.StakeAmount.ToBig(), *c.config.TransferBatchSubmissionGasLimit)
	arg1, arg2, arg3, arg4, arg5 := encoder.CommitmentsToTransferAndC2TSubmitBatchFields(batchID, commitments)
	tx, err := c.packAndRequest(&c.Rollup.Contract, opts, "submitTransfer", arg1, arg2, arg3, arg4, arg5)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.txsChannels.SentTxs <- tx

	return tx, nil
}

func (c *Client) SubmitCreate2TransfersBatch(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*types.Transaction, error) {
	opts := c.transactOpts(c.config.StakeAmount.ToBig(), *c.config.C2TBatchSubmissionGasLimit)
	arg1, arg2, arg3, arg4, arg5 := encoder.CommitmentsToTransferAndC2TSubmitBatchFields(batchID, commitments)
	tx, err := c.packAndRequest(&c.Rollup.Contract, opts, "submitCreate2Transfer", arg1, arg2, arg3, arg4, arg5)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.txsChannels.SentTxs <- tx

	return tx, nil
}

func (c *Client) SubmitMassMigrationsBatch(
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (*types.Transaction, error) {
	opts := c.transactOpts(c.config.StakeAmount.ToBig(), *c.config.MMBatchSubmissionGasLimit)
	arg1, arg2, arg3, arg4, arg5, arg6 := encoder.CommitmentsToSubmitMMBatchFields(batchID, commitments)
	tx, err := c.packAndRequest(&c.Rollup.Contract, opts, "submitMassMigration", arg1, arg2, arg3, arg4, arg5, arg6)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.txsChannels.SentTxs <- tx

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
	batch, err := c.GetBatch(models.NewUint256FromBig(*event.BatchID))
	if err != nil {
		return nil, err
	}
	batch.AccountTreeRoot = ref.Hash(common.BytesToHash(event.AccountRoot[:]))
	return batch, nil
}
