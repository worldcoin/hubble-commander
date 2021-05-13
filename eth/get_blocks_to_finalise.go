package eth

import "github.com/Worldcoin/hubble-commander/utils/ref"

func (c *Client) GetBlocksToFinalise() (*int64, error) {
	if c.blocksToFinalise != nil {
		return c.blocksToFinalise, nil
	}

	blocksToFinalise, err := c.rollup().ParamBlocksToFinalise()
	if err != nil {
		return nil, err
	}
	c.blocksToFinalise = ref.Int64(blocksToFinalise.Int64())
	return c.blocksToFinalise, nil
}
