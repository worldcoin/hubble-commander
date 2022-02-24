//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type DisputesE2ETestSuite struct {
	setup.E2ETestSuite
}

func (s *DisputesE2ETestSuite) SetupTest() {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = 32
	cfg.Rollup.MaxTxsPerCommitment = 32
	cfg.Rollup.MinCommitmentsPerBatch = 1

	s.SetupTestEnvironment(cfg)
}

func (s *DisputesE2ETestSuite) TestDisputes() {
	s.testDisputeSignatureTransfer()
	s.testDisputeSignatureC2T()
	s.testDisputeSignatureMM()

	s.testDisputeTransitionTransfer()
	s.testDisputeTransitionC2T()
	s.testDisputeTransitionMM()
}

func (s *DisputesE2ETestSuite) TestMeasureDisputesGasUsage() {
	s.measureDisputeSignatureTransfer()
	s.measureDisputeSignatureC2T()
	s.measureDisputeSignatureMM()

	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(0),
		})
	})

	s.measureDisputeTransitionTransfer()
	s.measureDisputeTransitionC2T()
	s.measureDisputeTransitionMM()
}

func TestDisputesE2ETestSuite(t *testing.T) {
	suite.Run(t, new(DisputesE2ETestSuite))
}

func (s *DisputesE2ETestSuite) testDisputeSignatureTransfer() {
	s.requireRollbackCompleted(func() {
		s.sendNTransfersBatchWithInvalidSignature(1, 1)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) testDisputeSignatureC2T() {
	s.requireRollbackCompleted(func() {
		s.sendNC2TsBatchWithInvalidSignature(1, 1)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) testDisputeSignatureMM() {
	s.requireRollbackCompleted(func() {
		s.sendNMMsBatchWithInvalidSignature(1, 1)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) testDisputeTransitionTransfer() {
	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(0),
		})
	})

	s.requireRollbackCompleted(func() {
		s.sendNTransfersBatchWithInvalidStateRoot(2, 1)
	})

	s.requireBatchesCount(2)

	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(32),
		})
	})
}

func (s *DisputesE2ETestSuite) testDisputeTransitionC2T() {
	s.requireRollbackCompleted(func() {
		s.sendNC2TsBatchWithInvalidStateRoot(3, 1)
	})

	s.requireBatchesCount(3)

	s.SubmitTxBatchAndWait(func() common.Hash {
		firstC2TWallet := s.Wallets[len(s.Wallets)-32]
		return s.SendNTransactions(32, dto.Create2Transfer{
			FromStateID: ref.Uint32(1),
			ToPublicKey: firstC2TWallet.PublicKey(),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(64),
		})
	})
}

func (s *DisputesE2ETestSuite) testDisputeTransitionMM() {
	s.requireRollbackCompleted(func() {
		s.sendNMMsBatchWithInvalidStateRoot(4, 1)
	})

	s.requireBatchesCount(4)

	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.MassMigration{
			FromStateID: ref.Uint32(1),
			SpokeID:     ref.Uint32(1),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(96),
		})
	})
}

func (s *DisputesE2ETestSuite) measureDisputeSignatureTransfer() {
	s.requireRollbackCompleted(func() {
		s.sendNTransfersBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) measureDisputeSignatureC2T() {
	s.requireRollbackCompleted(func() {
		s.sendNC2TsBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) measureDisputeSignatureMM() {
	s.requireRollbackCompleted(func() {
		s.sendNMMsBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) measureDisputeTransitionTransfer() {
	s.requireRollbackCompleted(func() {
		s.sendNTransfersBatchWithInvalidStateRoot(2, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(2)
}

func (s *DisputesE2ETestSuite) measureDisputeTransitionC2T() {
	s.requireRollbackCompleted(func() {
		s.sendNC2TsBatchWithInvalidStateRoot(2, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(2)
}

func (s *DisputesE2ETestSuite) measureDisputeTransitionMM() {
	s.requireRollbackCompleted(func() {
		s.sendNMMsBatchWithInvalidStateRoot(2, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(2)
}

func newEthClient(t *testing.T, client jsonrpc.RPCClient) *eth.Client {
	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	require.NoError(t, err)

	chainState := models.ChainState{
		ChainID:                        info.ChainID,
		AccountRegistry:                info.AccountRegistry,
		AccountRegistryDeploymentBlock: info.AccountRegistryDeploymentBlock,
		TokenRegistry:                  info.TokenRegistry,
		SpokeRegistry:                  info.SpokeRegistry,
		DepositManager:                 info.DepositManager,
		WithdrawManager:                info.WithdrawManager,
		Rollup:                         info.Rollup,
	}

	cfg := config.GetConfig()
	blockchain, err := chain.NewRPCConnection(cfg.Ethereum)
	require.NoError(t, err)

	backend := blockchain.GetBackend()

	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, backend)
	require.NoError(t, err)

	spokeRegistry, err := spokeregistry.NewSpokeRegistry(chainState.SpokeRegistry, backend)
	require.NoError(t, err)

	tokenRegistry, err := tokenregistry.NewTokenRegistry(chainState.TokenRegistry, backend)
	require.NoError(t, err)

	depositManager, err := depositmanager.NewDepositManager(chainState.DepositManager, backend)
	require.NoError(t, err)

	rollupContract, err := rollup.NewRollup(chainState.Rollup, backend)
	require.NoError(t, err)

	txsChannels := &eth.TxsTrackingChannels{
		SkipSentTxsChannel:                true,
		SkipSendingRequestsThroughChannel: true,
	}

	ethClient, err := eth.NewClient(blockchain, metrics.NewCommanderMetrics(), &eth.NewClientParams{
		ChainState:      chainState,
		AccountRegistry: accountRegistry,
		SpokeRegistry:   spokeRegistry,
		TokenRegistry:   tokenRegistry,
		DepositManager:  depositManager,
		Rollup:          rollupContract,
		TxsChannels:     txsChannels,
	})
	require.NoError(t, err)
	return ethClient
}
