// nolint: gocritic
package eth

import (
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
	"go.opentelemetry.io/otel/attribute"
)

type packAndRequestFunc func(shouldTrackTx bool, method string, data ...interface{}) (*types.Transaction, error)

type rollupSessionBuilder struct {
	rollup.RollupSession
	contract       Contract
	packAndRequest packAndRequestFunc

	attributes []attribute.KeyValue
	ctx        context.Context
}

func (c *Client) rollup() *rollupSessionBuilder {
	builder := rollupSessionBuilder{
		RollupSession: rollup.RollupSession{
			Contract:     c.Rollup.Rollup,
			TransactOpts: *c.Blockchain.GetAccount(),
		},
		contract:   c.Rollup.Contract,
		attributes: make([]attribute.KeyValue, 0),
		ctx:        context.Background(),
	}

	builder.packAndRequest = func(shouldTrackTx bool, method string, data ...interface{}) (*types.Transaction, error) {
		return c.packAndRequest(
			builder.ctx,
			&builder.contract,
			"Rollup",
			builder.attributes,
			&builder.TransactOpts,
			shouldTrackTx,
			method,
			data...,
		)
	}

	return &builder
}

func (b *rollupSessionBuilder) WithAttribute(kv attribute.KeyValue) *rollupSessionBuilder {
	b.attributes = append(b.attributes, kv)
	return b
}

func (b *rollupSessionBuilder) WithContext(ctx context.Context) *rollupSessionBuilder {
	b.ctx = ctx
	return b
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
	return b.packAndRequest(true, "withdrawStake", batchID)
}

func (b *rollupSessionBuilder) SubmitTransfer(
	batchID *big.Int,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	feeReceivers []*big.Int,
	txss [][]byte,
) (*types.Transaction, error) {
	return b.packAndRequest(true, "submitTransfer", batchID, stateRoots, signatures, feeReceivers, txss)
}

func (b *rollupSessionBuilder) SubmitCreate2Transfer(
	batchID *big.Int,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	feeReceivers []*big.Int,
	txss [][]byte,
) (*types.Transaction, error) {
	return b.packAndRequest(true, "submitCreate2Transfer", batchID, stateRoots, signatures, feeReceivers, txss)
}

func (b *rollupSessionBuilder) SubmitMassMigration(
	batchID *big.Int,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	meta [][4]*big.Int,
	withdrawRoots [][32]byte,
	txss [][]byte,
) (*types.Transaction, error) {
	return b.packAndRequest(true, "submitMassMigration", batchID, stateRoots, signatures, meta, withdrawRoots, txss)
}

func (b *rollupSessionBuilder) SubmitDeposits(
	batchID *big.Int,
	previous rollup.TypesCommitmentInclusionProof,
	vacant rollup.TypesSubtreeVacancyProof,
) (*types.Transaction, error) {
	return b.packAndRequest(true, "submitDeposits", batchID, previous, vacant)
}

func (b *rollupSessionBuilder) DisputeSignatureTransfer(
	batchID *big.Int,
	target rollup.TypesTransferCommitmentInclusionProof,
	signatureProof rollup.TypesSignatureProof,
) (*types.Transaction, error) {
	return b.packAndRequest(false, "disputeSignatureTransfer", batchID, target, signatureProof)
}

func (b *rollupSessionBuilder) DisputeSignatureCreate2Transfer(
	batchID *big.Int,
	target rollup.TypesTransferCommitmentInclusionProof,
	signatureProof rollup.TypesSignatureProofWithReceiver,
) (*types.Transaction, error) {
	return b.packAndRequest(false, "disputeSignatureCreate2Transfer", batchID, target, signatureProof)
}

func (b *rollupSessionBuilder) DisputeSignatureMassMigration(
	batchID *big.Int,
	target rollup.TypesMMCommitmentInclusionProof,
	signatureProof rollup.TypesSignatureProof,
) (*types.Transaction, error) {
	return b.packAndRequest(false, "disputeSignatureMassMigration", batchID, target, signatureProof)
}

func (b *rollupSessionBuilder) DisputeTransitionTransfer(
	batchID *big.Int,
	previous rollup.TypesCommitmentInclusionProof,
	target rollup.TypesTransferCommitmentInclusionProof,
	proofs []rollup.TypesStateMerkleProof,
) (*types.Transaction, error) {
	return b.packAndRequest(false, "disputeTransitionTransfer", batchID, previous, target, proofs)
}

func (b *rollupSessionBuilder) DisputeTransitionCreate2Transfer(
	batchID *big.Int,
	previous rollup.TypesCommitmentInclusionProof,
	target rollup.TypesTransferCommitmentInclusionProof,
	proofs []rollup.TypesStateMerkleProof,
) (*types.Transaction, error) {
	return b.packAndRequest(false, "disputeTransitionCreate2Transfer", batchID, previous, target, proofs)
}

func (b *rollupSessionBuilder) DisputeTransitionMassMigration(
	batchID *big.Int,
	previous rollup.TypesCommitmentInclusionProof,
	target rollup.TypesMMCommitmentInclusionProof,
	proofs []rollup.TypesStateMerkleProof,
) (*types.Transaction, error) {
	return b.packAndRequest(false, "disputeTransitionMassMigration", batchID, previous, target, proofs)
}

func (b *rollupSessionBuilder) KeepRollingBack() (*types.Transaction, error) {
	return b.packAndRequest(false, "keepRollingBack")
}

type accountRegistrySessionBuilder struct {
	accountregistry.AccountRegistrySession
	contract       Contract
	packAndRequest packAndRequestFunc

	attributes []attribute.KeyValue
	ctx        context.Context
}

func (a *AccountManager) accountRegistry() *accountRegistrySessionBuilder {
	builder := accountRegistrySessionBuilder{
		AccountRegistrySession: accountregistry.AccountRegistrySession{
			Contract:     a.AccountRegistry.AccountRegistry,
			TransactOpts: *a.Blockchain.GetAccount(),
		},
		contract:   a.AccountRegistry.Contract,
		attributes: make([]attribute.KeyValue, 0),
		ctx:        context.Background(),
	}

	builder.packAndRequest = func(shouldTrackTx bool, method string, data ...interface{}) (*types.Transaction, error) {
		return a.packAndRequest(
			builder.ctx,
			&builder.contract,
			builder.attributes,
			&builder.TransactOpts,
			shouldTrackTx,
			method,
			data...,
		)
	}

	return &builder
}

func (b *accountRegistrySessionBuilder) WithContext(ctx context.Context) *accountRegistrySessionBuilder {
	b.ctx = ctx
	return b
}

func (b *accountRegistrySessionBuilder) WithAttribute(kv attribute.KeyValue) *accountRegistrySessionBuilder {
	b.attributes = append(b.attributes, kv)
	return b
}

func (b *accountRegistrySessionBuilder) WithValue(value *models.Uint256) *accountRegistrySessionBuilder {
	b.TransactOpts.Value = value.ToBig()
	return b
}

func (b *accountRegistrySessionBuilder) WithGasLimit(gasLimit uint64) *accountRegistrySessionBuilder {
	b.TransactOpts.GasLimit = gasLimit
	return b
}

func (b *accountRegistrySessionBuilder) RegisterBatch(pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return b.packAndRequest(true, "registerBatch", pubkeys)
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
