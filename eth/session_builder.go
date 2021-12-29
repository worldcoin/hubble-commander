package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/models"
)

type rollupSessionBuilder struct {
	rollup.RollupSession
}

func (b *rollupSessionBuilder) WithValue(value *models.Uint256) *rollupSessionBuilder {
	b.TransactOpts.Value = value.ToBig()
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

func (b *accountRegistrySessionBuilder) WithValue(value *models.Uint256) *accountRegistrySessionBuilder {
	b.TransactOpts.Value = value.ToBig()
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
		Contract:     c.DepositManager.DepositManager,
		TransactOpts: *c.Blockchain.GetAccount(),
	}}
}

func (b *depositManagerSessionBuilder) WithValue(value *models.Uint256) *depositManagerSessionBuilder {
	b.TransactOpts.Value = value.ToBig()
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

func (b *tokenRegistrySessionBuilder) WithValue(value *models.Uint256) *tokenRegistrySessionBuilder {
	b.TransactOpts.Value = value.ToBig()
	return b
}

func (b *tokenRegistrySessionBuilder) WithGasLimit(gasLimit uint64) *tokenRegistrySessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}

type spokeRegistrySessionBuilder struct {
	spokeregistry.SpokeRegistrySession
}

func (c *Client) spokeRegistry() *spokeRegistrySessionBuilder {
	return &spokeRegistrySessionBuilder{spokeregistry.SpokeRegistrySession{
		Contract:     c.SpokeRegistry.SpokeRegistry,
		TransactOpts: *c.Blockchain.GetAccount(),
	}}
}

func (b *spokeRegistrySessionBuilder) WithValue(value *models.Uint256) *spokeRegistrySessionBuilder {
	b.TransactOpts.Value = value.ToBig()
	return b
}

func (b *spokeRegistrySessionBuilder) WithGasLimit(gasLimit uint64) *spokeRegistrySessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}
