//go:build e2e
// +build e2e

package bench

import (
	"fmt"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// TxTypeDistribution distribution of the transaction types sent by the test script
// example: { txtype.Create2Transfer: 0.2, txtype.Transfer: 0.8 } would mean 20% C2T, 80% transfers
type TxTypeDistribution = map[txtype.TransactionType]float32

type BenchmarkTransactionsSuite struct {
	benchmarkTestSuite
}

func (s *BenchmarkTransactionsSuite) SetupTest() {
	s.benchmarkTestSuite.SetupTest(BenchmarkConfig{
		TxAmount:               1_000,
		TxBatchSize:            32,
		MaxQueuedBatchesAmount: 20,
		MaxConcurrentWorkers:   4,
	})
}

func (s *BenchmarkTransactionsSuite) TestBenchTransfersCommander() {
	s.sendTransactions(TxTypeDistribution{txtype.Transfer: 1.0})
}

func (s *BenchmarkTransactionsSuite) TestBenchCreate2TransfersCommander() {
	s.sendTransactions(TxTypeDistribution{txtype.Create2Transfer: 1.0})
}

func (s *BenchmarkTransactionsSuite) TestBenchMixedCommander() {
	s.sendTransactions(TxTypeDistribution{txtype.Create2Transfer: 0.2, txtype.Transfer: 0.8}) // 20% C2T, 80% transfers
}

func (s *BenchmarkTransactionsSuite) TestBenchSyncCommander() {
	s.sendTransactions(TxTypeDistribution{txtype.Create2Transfer: 0.2, txtype.Transfer: 0.8})
	s.benchSyncing()
}

func (s *BenchmarkTransactionsSuite) benchSyncing() {
	cfg := config.GetConfig()

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5002"
	cfg.Badger.Path += "_passive"
	cfg.Bootstrap.ChainSpecPath = nil
	cfg.Bootstrap.BootstrapNodeURL = ref.String("http://localhost:8080")
	cfg.Ethereum.PrivateKey = "ab6919fd6ac00246bb78657e0696cf72058a4cb395133d074eabaddb83d8b00c"
	passiveCommander, err := setup.CreateInProcessCommanderWithConfig(cfg, false)
	s.NoError(err)
	err = passiveCommander.Start()
	s.NoError(err)
	defer func() {
		s.NoError(passiveCommander.Stop())
	}()

	// Observe commander syncing
	var networkInfo dto.NetworkInfo
	err = s.commander.Client().CallFor(&networkInfo, "hubble_getNetworkInfo")
	s.NoError(err)

	latestBatch := networkInfo.LatestBatch.Uint64()
	startTime := time.Now()
	lastSyncedBatch := uint64(0)
	for lastSyncedBatch < latestBatch {
		var networkInfo dto.NetworkInfo
		err = passiveCommander.Client().CallFor(&networkInfo, "hubble_getNetworkInfo")
		s.NoError(err)
		newBatch := uint64(0)
		if networkInfo.LatestBatch != nil {
			newBatch = networkInfo.LatestBatch.Uint64()
		}

		if newBatch == lastSyncedBatch {
			continue
		}
		lastSyncedBatch = newBatch

		txCount := networkInfo.TransactionCount

		fmt.Printf(
			"Transfers synced: %d, throughput: %f tx/s, batches synced: %d/%d\n",
			txCount,
			float64(txCount)/(time.Since(startTime).Seconds()),
			lastSyncedBatch,
			latestBatch,
		)
	}
}

func (s *BenchmarkTransactionsSuite) sendTransactions(distribution TxTypeDistribution) {
	s.prepareWorkers(func(senderWallet bls.Wallet, senderStateID uint32, nonce models.Uint256) common.Hash {
		var lastTxHash common.Hash

		txType := pickTxType(distribution)

		switch txType {
		case txtype.Transfer:
			// Pick random receiver id that's different from sender's.
			to := s.stateIds[randomInt(len(s.stateIds))]
			for to == senderStateID {
				to = s.stateIds[randomInt(len(s.stateIds))]
			}

			lastTxHash = s.sendTransfer(senderWallet, senderStateID, to, nonce)
		case txtype.Create2Transfer:
			// Pick random receiver pubkey
			to := s.wallets[randomInt(len(s.wallets))].PublicKey()

			lastTxHash = s.sendC2T(senderWallet, senderStateID, to, nonce)
		case txtype.MassMigration:
			panic("Not supported")
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

	log.Fatal("Unreachable")
	return txtype.Transfer
}

func TestBenchmarkTransactionsSuite(t *testing.T) {
	suite.Run(t, new(BenchmarkTransactionsSuite))
}
