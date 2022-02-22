package eth

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) WithdrawStakeAndWait(batchID *models.Uint256) error {
	tx, err := c.WithdrawStake(batchID)
	if err != nil {
		return err
	}
	_, err = c.WaitToBeMined(tx)
	return err
}

func (c *Client) WithdrawStake(batchID *models.Uint256) (*types.Transaction, error) {
	tx, err := c.rollup().
		WithGasLimit(*c.config.StakeWithdrawalGasLimit).
		WithdrawStake(batchID.ToBig())
	if err != nil {
		return nil, err
	}
	return tx, nil
}
