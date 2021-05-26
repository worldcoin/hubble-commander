package eth

import "github.com/Worldcoin/hubble-commander/contracts/chooser"

// IsActiveProposer checks if the current wallet is the active batch proposer.
func (c *Client) IsActiveProposer() (bool, error) {
	chooserAddress, err := c.Rollup.Chooser(nil)
	if err != nil {
		return false, err
	}

	chooserContract, err := chooser.NewChooser(chooserAddress, c.ChainConnection.GetBackend())
	if err != nil {
		return false, err
	}

	currentProposer, err := chooserContract.GetProposer(nil)
	if err != nil {
		return false, err
	}

	currentAddress := c.ChainConnection.GetAccount().From

	return currentAddress == currentProposer, nil
}
