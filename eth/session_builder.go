package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
)

type RollupSessionBuilder struct {
	rollup.RollupSession
}

func (b *RollupSessionBuilder) WithValue(value big.Int) *RollupSessionBuilder {
	b.TransactOpts.Value = &value
	return b
}

func (c *Client) rollup() *RollupSessionBuilder {
	return &RollupSessionBuilder{rollup.RollupSession{
		Contract:     c.Rollup,
		TransactOpts: *c.ChainConnection.GetAccount(),
	}}
}
