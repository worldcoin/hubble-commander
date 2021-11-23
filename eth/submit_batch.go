package eth

import (
	"context"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const gasEstimateMultiplier = 1.3

type SubmitBatchFunc func() (*types.Transaction, error)

func (c *Client) SubmitTransfersBatch(commitments []models.CommitmentWithTxs) (
	*types.Transaction,
	error,
) {
	input, err := c.packCommitments("submitTransfer", commitments)
	if err != nil {
		return nil, err
	}
	estimate, err := c.estimateBatchSubmissionGasLimit(input)
	if err != nil {
		return nil, err
	}
	return c.RawTransact(c.config.StakeAmount.ToBig(), estimate, input)
}

func (c *Client) SubmitCreate2TransfersBatch(commitments []models.CommitmentWithTxs) (
	*types.Transaction,
	error,
) {
	input, err := c.packCommitments("submitCreate2Transfer", commitments)
	if err != nil {
		return nil, err
	}
	estimate, err := c.estimateBatchSubmissionGasLimit(input)
	if err != nil {
		return nil, err
	}
	return c.RawTransact(c.config.StakeAmount.ToBig(), estimate, input)
}

func (c *Client) SubmitTransfersBatchAndWait(commitments []models.CommitmentWithTxs) (*models.Batch, error) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitTransfersBatch(commitments)
	})
}
func (c *Client) SubmitCreate2TransfersBatchAndWait(commitments []models.CommitmentWithTxs) (*models.Batch, error) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitCreate2TransfersBatch(commitments)
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

func (c *Client) estimateBatchSubmissionGasLimit(input []byte) (uint64, error) {
	account := c.Blockchain.GetAccount()
	msg := &ethereum.CallMsg{
		From:     account.From,
		To:       &c.ChainState.Rollup,
		GasPrice: account.GasPrice,
		Value:    c.config.StakeAmount.ToBig(),
		Data:     input,
	}
	estimatedGas, err := c.Blockchain.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, err
	}
	return uint64(float64(estimatedGas) * gasEstimateMultiplier), nil
}

func (c *Client) packCommitments(method string, commitments []models.CommitmentWithTxs) ([]byte, error) {
	stateRoots, signatures, feeReceivers, transactions := encoder.CommitmentToCalldataFields(commitments)
	return c.Rollup.ABI.Pack(method, stateRoots, signatures, feeReceivers, transactions)
}
