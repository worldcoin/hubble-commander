package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
)

type rollupSessionBuilder struct {
	rollup.RollupSession
	contract       Contract
	packAndRequest packAndRequestFunc
}

func (b *rollupSessionBuilder) WithValue(value *models.Uint256) *rollupSessionBuilder {
	b.TransactOpts.Value = value.ToBig()
	return b
}

func (b *rollupSessionBuilder) WithGasLimit(gasLimit uint64) *rollupSessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}

func (b *rollupSessionBuilder) WithdrawStake(batchID *big.Int) (*types.Transaction, error) {
	return b.packAndRequest(&b.contract, &b.TransactOpts, "withdrawStake", batchID)
}

func (b *rollupSessionBuilder) SubmitTransfer(
	batchID *big.Int,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	feeReceivers []*big.Int,
	txss [][]byte,
) (*types.Transaction, error) {
	return b.packAndRequest(&b.contract, &b.TransactOpts, "submitTransfer", batchID, stateRoots, signatures, feeReceivers, txss)
}

func (b *rollupSessionBuilder) SubmitCreate2Transfer(
	batchID *big.Int,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	feeReceivers []*big.Int,
	txss [][]byte,
) (*types.Transaction, error) {
	return b.packAndRequest(&b.contract, &b.TransactOpts, "submitCreate2Transfer", batchID, stateRoots, signatures, feeReceivers, txss)
}

func (b *rollupSessionBuilder) SubmitMassMigration(
	batchID *big.Int,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	meta [][4]*big.Int,
	withdrawRoots [][32]byte,
	txss [][]byte,
) (*types.Transaction, error) {
	return b.packAndRequest(&b.contract, &b.TransactOpts, "submitMassMigration", batchID, stateRoots, signatures, meta, withdrawRoots, txss)
}

func (b *rollupSessionBuilder) SubmitDeposits(
	batchID *big.Int,
	previous rollup.TypesCommitmentInclusionProof, //nolint:gocritic
	vacant rollup.TypesSubtreeVacancyProof,
) (*types.Transaction, error) {
	return b.packAndRequest(&b.contract, &b.TransactOpts, "submitDeposits", batchID, previous, vacant)
}

func (c *Client) rollup() *rollupSessionBuilder {
	return &rollupSessionBuilder{
		RollupSession: rollup.RollupSession{
			Contract:     c.Rollup.Rollup,
			TransactOpts: *c.Blockchain.GetAccount(),
		},
		contract:       c.Rollup.Contract,
		packAndRequest: c.packAndRequest,
	}
}

type accountRegistrySessionBuilder struct {
	accountregistry.AccountRegistrySession
	contract       Contract
	packAndRequest packAndRequestFunc
}

func (a *AccountManager) accountRegistry() *accountRegistrySessionBuilder {
	return &accountRegistrySessionBuilder{
		AccountRegistrySession: accountregistry.AccountRegistrySession{
			Contract:     a.AccountRegistry.AccountRegistry,
			TransactOpts: *a.Blockchain.GetAccount(),
		},
		contract:       a.AccountRegistry.Contract,
		packAndRequest: a.packAndRequest,
	}
}

func (b *accountRegistrySessionBuilder) WithValue(value *models.Uint256) *accountRegistrySessionBuilder {
	b.TransactOpts.Value = value.ToBig()
	return b
}

func (b *accountRegistrySessionBuilder) WithGasLimit(gasLimit uint64) *accountRegistrySessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}

//nolint: gocritic
func (b *accountRegistrySessionBuilder) RegisterBatch(pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return b.packAndRequest(&b.contract, &b.TransactOpts, "registerBatch", pubkeys)
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
