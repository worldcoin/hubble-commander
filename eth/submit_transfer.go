package eth

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SubmitTransferFunc func(commitments []models.Commitment) (*types.Transaction, error)

func (c *Client) SubmitTransfer() SubmitTransferFunc {
	return func(commitments []models.Commitment) (*types.Transaction, error) {
		return c.rollup().
			WithValue(c.config.stakeAmount.Int).
			SubmitCreate2Transfer(parseCommitments(commitments))
	}
}

func (c *Client) SubmitCreate2Transfer() SubmitTransferFunc {
	return func(commitments []models.Commitment) (*types.Transaction, error) {
		return c.rollup().
			WithValue(c.config.stakeAmount.Int).
			SubmitCreate2Transfer(parseCommitments(commitments))
	}
}

func (c *Client) SubmitTransfersBatch(commitments []models.Commitment, f SubmitTransferFunc) (batch *models.Batch, accountTreeRoot *common.Hash, err error) {
	sink := make(chan *rollup.RollupNewBatch)
	subscription, err := c.Rollup.WatchNewBatch(&bind.WatchOpts{}, sink)
	if err != nil {
		return
	}
	defer subscription.Unsubscribe()

	tx, err := f(commitments)
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
