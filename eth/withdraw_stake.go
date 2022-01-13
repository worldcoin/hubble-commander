package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) WithdrawStakeAndWait(batchID *models.Uint256) error {
	tx, err := c.rollup().WithdrawStake(batchID.ToBig())
	if err != nil {
		return err
	}
	_, err = chain.WaitToBeMined(c.Blockchain.GetBackend(), tx)
	return err
}

func (c *Client) WithdrawStake(batchID *models.Uint256) error {
	_, err := c.rollup().WithdrawStake(batchID.ToBig())
	return err
}

func (c *Client) RetrieveStakeWithdrawBatchID(receipt *types.Receipt) (*models.Uint256, error) {
	log, err := retrieveLog(receipt, StakeWithdrawEvent)
	if err != nil {
		return nil, err
	}

	event := new(rollup.RollupStakeWithdraw)
	err = c.SpokeRegistry.BoundContract.UnpackLog(event, StakeWithdrawEvent, *log)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return models.NewUint256FromBig(*event.BatchID), nil
}
