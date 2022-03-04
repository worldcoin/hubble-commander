// nolint: gocritic
package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
)

type packAndRequestFunc func(shouldTrackTx bool, method string, data ...interface{}) (*types.Transaction, error)

type sessionBuildersCreator interface {
	rollup() *rollupSessionBuilder
	depositManager() *depositManagerSessionBuilder
	tokenRegistry() *tokenRegistrySessionBuilder
	spokeRegistry() *spokeRegistrySessionBuilder
}

type accountRegistrySessionBuilderCreator interface {
	accountRegistry() *accountRegistrySessionBuilder
}

type clientSessionBuilders struct {
	blockchain             chain.Connection
	requestsChan           chan<- *TxSendingRequest
	rollupContract         *Rollup
	depositManagerContract *DepositManager
	tokenRegistryContract  *TokenRegistry
	spokeRegistryContract  *SpokeRegistry
}

type testClientSessionBuilders struct {
	*clientSessionBuilders
}

type accountManagerSessionBuilder struct {
	blockchain              chain.Connection
	requestsChan            chan<- *TxSendingRequest
	accountRegistryContract *AccountRegistry
}

type testAccountManagerSessionBuilder struct {
	*accountManagerSessionBuilder
}

func newSessionBuilders(
	blockchain chain.Connection,
	requestsChan chan<- *TxSendingRequest,
	rollupContract *Rollup,
	trContract *TokenRegistry,
	dmContract *DepositManager,
	srContract *SpokeRegistry,
) *clientSessionBuilders {
	return &clientSessionBuilders{
		blockchain:             blockchain,
		requestsChan:           requestsChan,
		rollupContract:         rollupContract,
		tokenRegistryContract:  trContract,
		depositManagerContract: dmContract,
		spokeRegistryContract:  srContract,
	}
}

func newTestSessionBuilders(
	blockchain chain.Connection,
	rollupContract *Rollup,
	dmContract *DepositManager,
	trContract *TokenRegistry,
	srContract *SpokeRegistry,
) *testClientSessionBuilders {
	return &testClientSessionBuilders{
		clientSessionBuilders: &clientSessionBuilders{
			blockchain:             blockchain,
			rollupContract:         rollupContract,
			depositManagerContract: dmContract,
			tokenRegistryContract:  trContract,
			spokeRegistryContract:  srContract,
		},
	}
}

func newAccountManagerSessionBuilder(
	blockchain chain.Connection,
	requestsChan chan<- *TxSendingRequest,
	arContract *AccountRegistry,
) *accountManagerSessionBuilder {
	return &accountManagerSessionBuilder{
		blockchain:              blockchain,
		requestsChan:            requestsChan,
		accountRegistryContract: arContract,
	}
}

func newTestAccountManagerSessionBuilder(
	blockchain chain.Connection,
	arContract *AccountRegistry,
) *testAccountManagerSessionBuilder {
	return &testAccountManagerSessionBuilder{
		accountManagerSessionBuilder: &accountManagerSessionBuilder{
			blockchain:              blockchain,
			accountRegistryContract: arContract,
		},
	}
}

type rollupSessionBuilder struct {
	rollup.RollupSession
	contract       Contract
	packAndRequest packAndRequestFunc
}

func (s *clientSessionBuilders) rollup() *rollupSessionBuilder {
	builder := rollupSessionBuilder{
		RollupSession: rollup.RollupSession{
			Contract:     s.rollupContract.Rollup,
			TransactOpts: *s.blockchain.GetAccount(),
		},
		contract: s.rollupContract.Contract,
	}

	builder.packAndRequest = func(shouldTrackTx bool, method string, data ...interface{}) (*types.Transaction, error) {
		return packAndRequest(s.requestsChan, &builder.contract, &builder.TransactOpts, shouldTrackTx, method, data...)
	}

	return &builder
}

func (s *testClientSessionBuilders) rollup() *rollupSessionBuilder {
	builder := rollupSessionBuilder{
		RollupSession: rollup.RollupSession{
			Contract:     s.rollupContract.Rollup,
			TransactOpts: *s.blockchain.GetAccount(),
		},
		contract: s.rollupContract.Contract,
	}

	builder.packAndRequest = func(shouldTrackTx bool, method string, data ...interface{}) (*types.Transaction, error) {
		return rawRequest(&builder.contract, &builder.TransactOpts, method, data...)
	}

	return &builder
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
}

func (a *accountManagerSessionBuilder) accountRegistry() *accountRegistrySessionBuilder {
	builder := accountRegistrySessionBuilder{
		AccountRegistrySession: accountregistry.AccountRegistrySession{
			Contract:     a.accountRegistryContract.AccountRegistry,
			TransactOpts: *a.blockchain.GetAccount(),
		},
		contract: a.accountRegistryContract.Contract,
	}

	builder.packAndRequest = func(shouldTrackTx bool, method string, data ...interface{}) (*types.Transaction, error) {
		return packAndRequest(a.requestsChan, &builder.contract, &builder.TransactOpts, shouldTrackTx, method, data...)
	}

	return &builder
}

func (a *testAccountManagerSessionBuilder) accountRegistry() *accountRegistrySessionBuilder {
	builder := accountRegistrySessionBuilder{
		AccountRegistrySession: accountregistry.AccountRegistrySession{
			Contract:     a.accountRegistryContract.AccountRegistry,
			TransactOpts: *a.blockchain.GetAccount(),
		},
		contract: a.accountRegistryContract.Contract,
	}

	builder.packAndRequest = func(shouldTrackTx bool, method string, data ...interface{}) (*types.Transaction, error) {
		return rawRequest(&builder.contract, &builder.TransactOpts, method, data...)
	}

	return &builder
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

func (s *clientSessionBuilders) depositManager() *depositManagerSessionBuilder {
	return &depositManagerSessionBuilder{depositmanager.DepositManagerSession{
		Contract:     s.depositManagerContract.DepositManager,
		TransactOpts: *s.blockchain.GetAccount(),
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

func (s *clientSessionBuilders) tokenRegistry() *tokenRegistrySessionBuilder {
	return &tokenRegistrySessionBuilder{tokenregistry.TokenRegistrySession{
		Contract:     s.tokenRegistryContract.TokenRegistry,
		TransactOpts: *s.blockchain.GetAccount(),
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

func (s *clientSessionBuilders) spokeRegistry() *spokeRegistrySessionBuilder {
	return &spokeRegistrySessionBuilder{spokeregistry.SpokeRegistrySession{
		Contract:     s.spokeRegistryContract.SpokeRegistry,
		TransactOpts: *s.blockchain.GetAccount(),
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
