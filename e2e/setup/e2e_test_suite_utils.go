package setup

import (
	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	// "github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/ethereum/go-ethereum/common"
)

func (s *E2ETestSuite) CalculateDepositsCountForFullBatch() int {
	subtreeDepth, err := s.ETHClient.GetMaxSubtreeDepthParam()
	s.NoError(err)
	depositsCount := 1 << *subtreeDepth

	return depositsCount
}

func (s *E2ETestSuite) getDomain() bls.Domain {
	info := s.GetNetworkInfo()
	return info.SignatureDomain
}

func (s *E2ETestSuite) newCommanderETHClient() *eth.Client {
	return s.newEthClient(s.Cfg.Ethereum.PrivateKey)
}

func (s *E2ETestSuite) newE2ETestETHClient() *eth.Client {
	return s.newEthClient(TestEthClientPrivateKey)
}

func (s *E2ETestSuite) newEthClient(privateKey string) *eth.Client {
	chainSpec := s.Commander.ChainSpec()
	chainState := models.ChainState{
		ChainID:                        chainSpec.ChainID,
		AccountRegistry:                chainSpec.AccountRegistry,
		AccountRegistryDeploymentBlock: chainSpec.AccountRegistryDeploymentBlock,
		TokenRegistry:                  chainSpec.TokenRegistry,
		SpokeRegistry:                  chainSpec.SpokeRegistry,
		DepositManager:                 chainSpec.DepositManager,
		WithdrawManager:                chainSpec.WithdrawManager,
		Rollup:                         chainSpec.Rollup,
		GenesisAccounts:                chainSpec.GenesisAccounts,
	}

	cfg := config.GetCommanderConfigAndSetupLogger()
	cfg.Ethereum.PrivateKey = privateKey
	blockchain, err := chain.NewRPCConnection(cfg.Ethereum)
	s.NoError(err)

	backend := blockchain.GetBackend()

	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, backend)
	s.NoError(err)

	spokeRegistry, err := spokeregistry.NewSpokeRegistry(chainState.SpokeRegistry, backend)
	s.NoError(err)

	tokenRegistry, err := tokenregistry.NewTokenRegistry(chainState.TokenRegistry, backend)
	s.NoError(err)

	depositManager, err := depositmanager.NewDepositManager(chainState.DepositManager, backend)
	s.NoError(err)

	rollupContract, err := rollup.NewRollup(chainState.Rollup, backend)
	s.NoError(err)

	txsChannels := &eth.TxsTrackingChannels{
		SkipChannelSending: true,
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
	s.NoError(err)
	return ethClient
}

func (s *E2ETestSuite) submitTransactionToCommander(tx interface{}) common.Hash {
	var txHash common.Hash
	err := s.RPCClient.CallFor(&txHash, "hubble_sendTransaction", []interface{}{tx})
	s.NoError(err)
	s.NotZero(txHash)
	return txHash
}

func (s *E2ETestSuite) sendTransfer(transfer dto.Transfer) common.Hash {
	signedTransfer, err := api.SignTransfer(&s.Wallets[*transfer.FromStateID], transfer)
	s.NoError(err)

	return s.submitTransactionToCommander(*signedTransfer)
}

func (s *E2ETestSuite) sendCreate2Transfer(transfer dto.Create2Transfer) common.Hash {
	signedTransfer, err := api.SignCreate2Transfer(&s.Wallets[*transfer.FromStateID], transfer)
	s.NoError(err)

	return s.submitTransactionToCommander(*signedTransfer)
}

func (s *E2ETestSuite) sendMassMigration(massMigration dto.MassMigration) common.Hash {
	signedMassMigration, err := api.SignMassMigration(&s.Wallets[*massMigration.FromStateID], massMigration)
	s.NoError(err)

	return s.submitTransactionToCommander(*signedMassMigration)
}

func (s *E2ETestSuite) sendNTransfers(n int, transfer dto.Transfer) common.Hash {
	firstTxHash := s.sendTransfer(transfer)
	// TODO fix this and add it back in
	/*
	firstTxReceipt := s.GetTransaction(firstTxHash)
	s.Equal(txstatus.Pending, firstTxReceipt.Status)
	*/
	for i := 1; i < n; i++ {
		transfer.Nonce = transfer.Nonce.AddN(1)
		s.sendTransfer(transfer)
	}

	return firstTxHash
}

func (s *E2ETestSuite) sendNCreate2Transfers(n int, transfer dto.Create2Transfer) common.Hash {
	firstTxHash := s.sendCreate2Transfer(transfer)
	/*
	firstTxReceipt := s.GetTransaction(firstTxHash)
	s.Equal(txstatus.Pending, firstTxReceipt.Status)
	*/
	for i := 1; i < n; i++ {
		transfer.Nonce = transfer.Nonce.AddN(1)
		s.sendCreate2Transfer(transfer)
	}

	return firstTxHash
}

func (s *E2ETestSuite) sendNMassMigrations(n int, massMigration dto.MassMigration) common.Hash {
	firstTxHash := s.sendMassMigration(massMigration)
	/*
	firstTxReceipt := s.GetTransaction(firstTxHash)
	s.Equal(txstatus.Pending, firstTxReceipt.Status)
	*/

	for i := 1; i < n; i++ {
		massMigration.Nonce = massMigration.Nonce.AddN(1)
		s.sendMassMigration(massMigration)
	}

	return firstTxHash
}
