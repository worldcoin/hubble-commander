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
	opts := *c.Blockchain.GetAccount()
	opts.GasLimit = *c.config.StakeWithdrawalGasLimit
	tx, err := c.packAndRequest(&c.Rollup.Contract, &opts, "withdrawStake", batchID.ToBig())
	if err != nil {
		return nil, err
	}
	return tx, nil
}
