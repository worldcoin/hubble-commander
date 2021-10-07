package eth

import (
	"context"
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const gasEstimateMultiplier = 1.3

var ErrSubmitBatchAndWait = fmt.Errorf("submitBatchAndWait: timeout")

type SubmitBatchFunc func() (*types.Transaction, error)

func (c *Client) SubmitTransfersBatch(commitments []models.Commitment) (
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

func (c *Client) SubmitCreate2TransfersBatch(commitments []models.Commitment) (
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

func (c *Client) SubmitTransfersBatchAndWait(commitments []models.Commitment) (*models.Batch, error) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitTransfersBatch(commitments)
	})
}
func (c *Client) SubmitCreate2TransfersBatchAndWait(commitments []models.Commitment) (*models.Batch, error) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitCreate2TransfersBatch(commitments)
	})
}

func (c *Client) submitBatchAndWait(submit SubmitBatchFunc) (batch *models.Batch, err error) {
	sink := make(chan *rollup.RollupNewBatch)
	subscription, err := c.Rollup.WatchNewBatch(&bind.WatchOpts{}, sink)
	if err != nil {
		return
	}
	defer subscription.Unsubscribe()

	tx, err := submit()
	if err != nil {
		return
	}

	for {
		select {
		case newBatch := <-sink:
			if newBatch.Raw.TxHash == tx.Hash() {
				return c.handleNewBatchEvent(newBatch)
			}
		case <-time.After(*c.config.TxTimeout):
			return nil, errors.WithStack(ErrSubmitBatchAndWait)
		}
	}
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

func (c *Client) packCommitments(method string, commitments []models.Commitment) ([]byte, error) {
	stateRoots, signatures, feeReceivers, transactions := encoder.CommitmentToCalldataFields(commitments)
	return c.RollupABI.Pack(method, stateRoots, signatures, feeReceivers, transactions)
}
