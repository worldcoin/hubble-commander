//go:build e2e
// +build e2e

package bench

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type BenchmarkPubKeyRegistrationSuite struct {
	benchmarkTestSuite

	unregisteredWallets []bls.Wallet
}

func (s *BenchmarkPubKeyRegistrationSuite) SetupTest() {
	s.benchmarkTestSuite.SetupTest(BenchmarkConfig{
		TxAmount:               1_000,
		TxBatchSize:            32,
		MaxQueuedBatchesAmount: 20,
		MaxConcurrentWorkers:   4,
	})

	unregisteredWallets, err := setup.CreateUnregisteredWalletsForBenchmark(s.benchConfig.TxAmount, s.domain)
	s.NoError(err)

	s.unregisteredWallets = unregisteredWallets
}

func (s *BenchmarkPubKeyRegistrationSuite) TestBenchPubKeysRegistration() {
	unregisteredWalletsIndex := 0
	s.sendTransactions(func(senderWallet bls.Wallet, senderStateID uint32, nonce models.Uint256) common.Hash {
		to := s.unregisteredWallets[unregisteredWalletsIndex].PublicKey()
		unregisteredWalletsIndex++

		return s.sendC2T(senderWallet, senderStateID, to, nonce)
	})
}

func TestBenchmarkPubKeyRegistrationSuite(t *testing.T) {
	suite.Run(t, new(BenchmarkPubKeyRegistrationSuite))
}
