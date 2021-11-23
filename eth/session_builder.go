package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
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
		Contract:     c.Rollup.Rollup,
		TransactOpts: *c.Blockchain.GetAccount(),
	}}
}

type accountRegistrySessionBuilder struct {
	accountregistry.AccountRegistrySession
}

func (a *AccountManager) accountRegistry() *accountRegistrySessionBuilder {
	return &accountRegistrySessionBuilder{accountregistry.AccountRegistrySession{
		Contract:     a.AccountRegistry.AccountRegistry,
		TransactOpts: *a.Blockchain.GetAccount(),
	}}
}

func (b *accountRegistrySessionBuilder) WithValue(value big.Int) *accountRegistrySessionBuilder {
	b.TransactOpts.Value = &value
	return b
}

func (b *accountRegistrySessionBuilder) WithGasLimit(gasLimit uint64) *accountRegistrySessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}

type depositManagerSessionBuilder struct {
	depositmanager.DepositManagerSession
}

func (c *Client) depositManager() *depositManagerSessionBuilder {
	return &depositManagerSessionBuilder{depositmanager.DepositManagerSession{
		Contract:     c.DepositManager,
		TransactOpts: *c.Blockchain.GetAccount(),
	}}
}

func (b *depositManagerSessionBuilder) WithValue(value big.Int) *depositManagerSessionBuilder {
	b.TransactOpts.Value = &value
	return b
}

func (b *depositManagerSessionBuilder) WithGasLimit(gasLimit uint64) *depositManagerSessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}

type tokenRegistrySessionBuilder struct {
	tokenregistry.TokenRegistrySession
}

func (c *Client) tokenRegistry() *tokenRegistrySessionBuilder {
	return &tokenRegistrySessionBuilder{tokenregistry.TokenRegistrySession{
		Contract:     c.TokenRegistry.TokenRegistry,
		TransactOpts: *c.Blockchain.GetAccount(),
	}}
}

func (b *tokenRegistrySessionBuilder) WithValue(value big.Int) *tokenRegistrySessionBuilder {
	b.TransactOpts.Value = &value
	return b
}

func (b *tokenRegistrySessionBuilder) WithGasLimit(gasLimit uint64) *tokenRegistrySessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}
