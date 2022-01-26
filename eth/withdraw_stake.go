package eth

import (
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
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
