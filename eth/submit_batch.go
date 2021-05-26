package eth

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SubmitBatchFunc func(commitments []models.Commitment) (*types.Transaction, error)

func (c *Client) SubmitTransfersBatchAndMine(commitments []models.Commitment) (
	batch *models.Batch,
	accountTreeRoot *common.Hash,
	err error,
) {
	return c.submitBatch(commitments, func(commitments []models.Commitment) (*types.Transaction, error) {
		return c.rollup().
			WithValue(*c.config.stakeAmount.ToBig()).
			SubmitTransfer(encoder.CommitmentToCalldataFields(commitments))
	})
}

func (c *Client) SubmitCreate2TransfersBatchAndMine(commitments []models.Commitment) (
	batch *models.Batch,
	accountTreeRoot *common.Hash,
	err error,
) {
	return c.submitBatch(commitments, func(commitments []models.Commitment) (*types.Transaction, error) {
		return c.rollup().
			WithValue(*c.config.stakeAmount.ToBig()).
			SubmitCreate2Transfer(encoder.CommitmentToCalldataFields(commitments))
	})
}

func (c *Client) SubmitTransfersBatch(commitments []models.Commitment) (
	*types.Transaction,
	error,
) {
	return c.rollup().
		WithValue(*c.config.stakeAmount.ToBig()).
		SubmitTransfer(encoder.CommitmentToCalldataFields(commitments))
}
func (c *Client) SubmitCreate2TransfersBatch(commitments []models.Commitment) (
	*types.Transaction,
	error,
) {
	return c.rollup().
		WithValue(*c.config.stakeAmount.ToBig()).
		SubmitCreate2Transfer(encoder.CommitmentToCalldataFields(commitments))
}

func (c *Client) submitBatch(
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
