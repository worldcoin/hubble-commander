package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/chooser"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// IsActiveProposer checks if the current wallet is the active batch proposer.
func (c *Client) IsActiveProposer() (bool, error) {
	chooserAddress, err := c.Rollup.Chooser(nil)
	if err != nil {
		return false, err
	}

	chooserContract, err := chooser.NewChooser(chooserAddress, c.Blockchain.GetBackend())
	if err != nil {
		return false, err
	}

	currentProposer, err := chooserContract.GetProposer(&bind.CallOpts{
		From: c.Blockchain.GetAccount().From,
	})
	if err != nil {
		return false, err
	}

	currentAddress := c.Blockchain.GetAccount().From

	return currentAddress == currentProposer, nil
}
