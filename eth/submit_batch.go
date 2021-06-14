package eth

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const gasEstimationFactor = 1.5

type SubmitBatchFunc func(commitments []models.Commitment) (*types.Transaction, error)

func (c *Client) SubmitTransfersBatch(commitments []models.Commitment) (
	*types.Transaction,
	error,
) {
	stateRoots, signatures, feeReceivers, transactions := encoder.CommitmentToCalldataFields(commitments)
	estimate, err := c.estimateGasLimit("submitTransfer", stateRoots, signatures, feeReceivers, transactions)
	if err != nil {
		return nil, err
	}
	return c.rollup().
		WithValue(*c.config.stakeAmount.ToBig()).
		WithGasLimit(estimate).
		SubmitTransfer(stateRoots, signatures, feeReceivers, transactions)
}

func (c *Client) SubmitCreate2TransfersBatch(commitments []models.Commitment) (
	*types.Transaction,
	error,
) {
	stateRoots, signatures, feeReceivers, transactions := encoder.CommitmentToCalldataFields(commitments)
	estimate, err := c.estimateGasLimit("submitCreate2Transfer", stateRoots, signatures, feeReceivers, transactions)
	if err != nil {
		return nil, err
	}
	return c.rollup().
		WithValue(*c.config.stakeAmount.ToBig()).
		WithGasLimit(estimate).
		SubmitCreate2Transfer(stateRoots, signatures, feeReceivers, transactions)
}

func (c *Client) SubmitTransfersBatchAndMine(commitments []models.Commitment) (
	batch *models.Batch,
	accountTreeRoot *common.Hash,
	err error,
) {
	return c.submitBatchAndMine(commitments, c.SubmitTransfersBatch)
}
func (c *Client) SubmitCreate2TransfersBatchAndMine(commitments []models.Commitment) (
	batch *models.Batch,
	accountTreeRoot *common.Hash,
	err error,
) {
	return c.submitBatchAndMine(commitments, c.SubmitCreate2TransfersBatch)
}

func (c *Client) submitBatchAndMine(
	commitments []models.Commitment,
	submit SubmitBatchFunc,
) (batch *models.Batch, accountTreeRoot *common.Hash, err error) {
	sink := make(chan *rollup.RollupNewBatch)
	subscription, err := c.Rollup.WatchNewBatch(&bind.WatchOpts{}, sink)
	if err != nil {
		return
	}
	defer subscription.Unsubscribe()

	tx, err := submit(commitments)
	if err != nil {
		return
	}

	for {
		select {
		case newBatch := <-sink:
			if newBatch.Raw.TxHash == tx.Hash() {
				return c.handleNewBatchEvent(newBatch)
			}
		case <-time.After(*c.config.txTimeout):
			return nil, nil, fmt.Errorf("timeout")
		}
	}
}

func (c *Client) handleNewBatchEvent(event *rollup.RollupNewBatch) (*models.Batch, *common.Hash, error) {
	batch, err := c.GetBatch(models.NewUint256FromBig(*event.BatchID))
	if err != nil {
		return nil, nil, err
	}
	accountRoot := common.BytesToHash(event.AccountRoot[:])
	return batch, &accountRoot, nil
}

func (c *Client) estimateGasLimit(
	method string,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	feeReceivers []*big.Int,
	transactions [][]byte) (uint64, error) {
	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	if err != nil {
		return 0, errors.WithStack(err)
	}
	input, err := rollupAbi.Pack(method, stateRoots, signatures, feeReceivers, transactions)
	if err != nil {
		return 0, err
	}
	account := c.ChainConnection.GetAccount()
	estimatedGas, err := c.ChainConnection.EstimateGas(context.Background(), &ethereum.CallMsg{
		From:     account.From,
		To:       &c.ChainState.Rollup,
		GasPrice: account.GasPrice,
		Value:    c.config.stakeAmount.ToBig(),
		Data:     input,
	})
	if err != nil {
		return 0, err
	}
	return uint64(float64(estimatedGas) * gasEstimationFactor), nil
}
