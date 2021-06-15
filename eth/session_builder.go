package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
)

type rollupSessionBuilder struct {
	rollup.RollupSession
}

func (b *rollupSessionBuilder) WithValue(value big.Int) *rollupSessionBuilder {
	b.TransactOpts.Value = &value
	return b
}

func (b *rollupSessionBuilder) WithGasLimit(gasLimit uint64) *rollupSessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}

func (c *Client) rollup() *rollupSessionBuilder {
	return &rollupSessionBuilder{rollup.RollupSession{
		Contract:     c.Rollup,
		TransactOpts: *c.ChainConnection.GetAccount(),
	}}
}
