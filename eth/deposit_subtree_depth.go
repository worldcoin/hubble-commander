package eth

import (
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Client) GetDepositSubtreeDepth() (*uint8, error) {
	if c.depositSubtreeDepth != nil {
		return c.depositSubtreeDepth, nil
	}

	depositSubtreeDepth, err := c.DepositManager.ParamMaxSubtreeDepth(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	c.depositSubtreeDepth = ref.Uint8(uint8(depositSubtreeDepth.Uint64()))
	return c.depositSubtreeDepth, nil
}
