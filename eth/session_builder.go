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
