//go:build e2e
// +build e2e

package bench

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

// TxTypeDistribution distribution of the transaction types sent by the test script
// example: { txtype.Create2Transfer: 0.2, txtype.Transfer: 0.8 } would mean 20% C2T, 80% transfers
type TxTypeDistribution = map[txtype.TransactionType]float32

type BenchmarkTransactionsSuite struct {
	benchmarkTestSuite

	unregisteredWallets []bls.Wallet
}

func (s *BenchmarkTransactionsSuite) SetupTest() {
	s.benchmarkTestSuite.SetupTest(BenchmarkConfig{
		TxCount:               1_000,
		TxBatchSize:           32,
		MaxQueuedBatchesCount: 20,
		MaxConcurrentWorkers:  4,
	})

	unregisteredWallets, err := setup.CreateUnregisteredWalletsForBenchmark(s.benchConfig.TxCount, s.Domain)
	s.NoError(err)

	s.unregisteredWallets = unregisteredWallets
}

func (s *BenchmarkTransactionsSuite) TestBenchTransfersCommander() {
	s.sendTransactionsWithDistribution(TxTypeDistribution{txtype.Transfer: 1.0})
}

func (s *BenchmarkTransactionsSuite) TestBenchCreate2TransfersCommander() {
	s.sendTransactionsWithDistribution(TxTypeDistribution{txtype.Create2Transfer: 1.0})
}

func (s *BenchmarkTransactionsSuite) TestBenchMassMigrationsCommander() {
	time.Sleep(time.Second * 15) // give the NewBlockLoop a change to notice the spoke
	s.sendTransactionsWithDistribution(TxTypeDistribution{txtype.MassMigration: 1.0})
}

func (s *BenchmarkTransactionsSuite) TestBenchMixedCommander() {
	time.Sleep(time.Second * 15) // give the NewBlockLoop a change to notice the spoke
	s.sendTransactionsWithDistribution(TxTypeDistribution{
		txtype.Transfer:        0.75,
		txtype.Create2Transfer: 0.2,
		txtype.MassMigration:   0.05,
	}) // 75% transfers, 20% C2T, 5% MM
}

func (s *BenchmarkTransactionsSuite) TestBenchSyncCommander() {
	s.sendTransactionsWithDistribution(TxTypeDistribution{
		txtype.Transfer:        0.75,
		txtype.Create2Transfer: 0.2,
		txtype.MassMigration:   0.05,
	})
	s.benchSyncing()
}

func (s *BenchmarkTransactionsSuite) benchSyncing() {
	passiveCommander := s.preparePassiveCommander()
	err := passiveCommander.Start()
	s.NoError(err)
	defer func() {
		s.NoError(passiveCommander.Stop())
	}()

	// Observe commander syncing
	networkInfo := s.GetNetworkInfo()

	// Further calls are only done to the passive commander
	s.RPCClient = passiveCommander.Client()

	latestBatch := networkInfo.LatestBatch.Uint64()
	startTime := time.Now()
	lastSyncedBatch := uint64(0)
	for lastSyncedBatch < latestBatch {
		passiveCommanderNetworkInfo := s.GetNetworkInfo()
		newBatch := uint64(0)
		if passiveCommanderNetworkInfo.LatestBatch != nil {
			newBatch = passiveCommanderNetworkInfo.LatestBatch.Uint64()
		}

		if newBatch == lastSyncedBatch {
			continue
		}
		lastSyncedBatch = newBatch

		txCount := passiveCommanderNetworkInfo.TransactionCount

		fmt.Printf(
			"Transfers synced: %d, throughput: %f tx/s, batches synced: %d/%d\n",
			txCount,
			float64(txCount)/(time.Since(startTime).Seconds()),
			lastSyncedBatch,
			latestBatch,
		)
	}
}

func (s *BenchmarkTransactionsSuite) preparePassiveCommander() *setup.InProcessCommander {
	newCommanderCfg := *s.Cfg
	newCommanderCfg.Bootstrap.Prune = true
	newCommanderCfg.API.Port = "5555"
	newCommanderCfg.Metrics.Port = "2222"
	newCommanderCfg.Badger.Path += "_bench_passive"
	newCommanderCfg.Bootstrap.BootstrapNodeURL = ref.String("http://localhost:8080")
	newCommanderCfg.Ethereum.PrivateKey = "ab6919fd6ac00246bb78657e0696cf72058a4cb395133d074eabaddb83d8b00c"
	passiveCommander, err := setup.CreateInProcessCommander(&newCommanderCfg, nil)
	s.NoError(err)

	return passiveCommander
}

func (s *BenchmarkTransactionsSuite) sendTransactionsWithDistribution(distribution TxTypeDistribution) {
	unregisteredWalletsIndex := 0

	s.sendTransactions(func(senderWallet bls.Wallet, senderStateID uint32, nonce models.Uint256) common.Hash {
		var lastTxHash common.Hash

		txType := pickTxType(distribution)

		switch txType {
		case txtype.Transfer:
			// Pick random receiver id that's different from sender's.
			to := s.stateIds[randomInt(len(s.stateIds))]
			for to == senderStateID {
				to = s.stateIds[randomInt(len(s.stateIds))]
			}

			lastTxHash = s.sendTransferFromWallet(senderWallet, senderStateID, to, nonce)
		case txtype.Create2Transfer:
			// Pick random unregistered receiver pubkey
			to := s.unregisteredWallets[unregisteredWalletsIndex].PublicKey()
			unregisteredWalletsIndex++

			lastTxHash = s.sendC2TFromWallet(senderWallet, senderStateID, to, nonce)
		case txtype.MassMigration:
			lastTxHash = s.sendMMFromWallet(senderWallet, senderStateID, nonce)
		}

		return lastTxHash
	})
}

// pickTxType picks a random transaction type based on the weighted distribution
func pickTxType(distribution TxTypeDistribution) txtype.TransactionType {
	sum := float32(0)
	for _, weight := range distribution {
		sum += weight
	}

	pick := randomFloat32() * sum

	for txType, weight := range distribution {
		if weight >= pick {
			return txType
		} else {
			pick -= weight
		}
	}

	panic("unreachable")
}

func randomInt(n int) int {
	//nolint:gosec
	return rand.Intn(n)
}

func randomFloat32() float32 {
	//nolint:gosec
	return rand.Float32()
}

func TestBenchmarkTransactionsSuite(t *testing.T) {
	suite.Run(t, new(BenchmarkTransactionsSuite))
}
