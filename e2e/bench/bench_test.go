//go:build e2e
// +build e2e

package bench

import (
	"sync"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/e2e"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Total number of transactions to be sent.
const txCount = 1_000

// Number of transaction that will be sent in a single batch (unrelated to rollup "batches").
const txBatchSize = 32

// Maximum number of tx batches in queue.
const maxQueuedBatchesCount = 20

// Maximum number of workers that send transactions.
const maxConcurrentWorkers = 4

// TxTypeDistribution distribution of the transaction types sent by the test script
// example: { txtype.Create2Transfer: 0.2, txtype.Transfer: 0.8 } would mean 20% C2T, 80% transfers
type TxTypeDistribution = map[txtype.TransactionType]float32

type BenchmarkSuite struct {
	*require.Assertions
	suite.Suite

	commander                setup.Commander
	wallets                  []bls.Wallet
	unregisteredWallets      []bls.Wallet
	unregisteredWalletsIndex int
	stateIds                 []uint32

	startTime time.Time
	waitGroup sync.WaitGroup

	// Only use atomic operations to increment those two counters.
	txsSent             int64
	txsQueued           int64
	lastReportedTxCount int64
}

func (s *BenchmarkSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *BenchmarkSuite) SetupTest() {
	commander, err := setup.NewCommanderFromEnv()
	s.NoError(err)

	err = commander.Start()
	s.NoError(err)

	domain := e2e.GetDomain(s.T(), commander.Client())
	wallets, err := setup.CreateWallets(domain)
	s.NoError(err)
	unregisteredWallets, err := setup.CreateUnregisteredWalletsForBenchmark(txCount, domain)
	s.NoError(err)

	s.commander = commander
	s.wallets = wallets
	s.unregisteredWallets = unregisteredWallets
	s.unregisteredWalletsIndex = 0
	s.stateIds = make([]uint32, 0)
	s.waitGroup = sync.WaitGroup{}
	s.txsSent = 0
	s.txsQueued = 0
	s.lastReportedTxCount = 0
}

func (s *BenchmarkSuite) TearDownTest() {
	s.NoError(s.commander.Stop())
}

func (s *BenchmarkSuite) TestBenchTransfersCommander() {
	s.sendTransactions(TxTypeDistribution{txtype.Transfer: 1.0})
}

func (s *BenchmarkSuite) TestBenchCreate2TransfersCommander() {
	s.sendTransactions(TxTypeDistribution{txtype.Create2Transfer: 1.0})
}

func (s *BenchmarkSuite) TestBenchMixedCommander() {
	s.sendTransactions(TxTypeDistribution{txtype.Create2Transfer: 0.2, txtype.Transfer: 0.8}) // 20% C2T, 80% transfers
}

func (s *BenchmarkSuite) TestBenchSyncCommander() {
	s.sendTransactions(TxTypeDistribution{txtype.Create2Transfer: 0.2, txtype.Transfer: 0.8})
	s.benchSyncing()
}

func (s *BenchmarkSuite) TestBenchPubKeysRegistration() {
	s.registerPubKeys()
}

func TestBenchmarkSuite(t *testing.T) {
	suite.Run(t, new(BenchmarkSuite))
}
