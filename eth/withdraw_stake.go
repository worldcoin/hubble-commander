package eth

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (c *Client) WithdrawStakeAndWait(batchID *models.Uint256) error {
	tx, err := c.rollup().WithdrawStake(batchID.ToBig())
	if err != nil {
		return err
	}
	_, err = c.WaitToBeMined(tx)
	return err
}

func (c *Client) WithdrawStake(batchID *models.Uint256) error {
	tx, err := c.rollup().
		WithGasLimit(200_000).
		WithdrawStake(batchID.ToBig())
	if err != nil {
		return err
	}
	c.txsHashesChan <- tx.Hash()
	return nil
}
