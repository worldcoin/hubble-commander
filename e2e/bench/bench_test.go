//go:build e2e
// +build e2e

package bench

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
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

	commander setup.Commander
	wallets   []bls.Wallet
	stateIds  []uint32

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

	s.commander = commander
	s.wallets = wallets
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

func (s *BenchmarkSuite) benchSyncing() {
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

func (s *BenchmarkSuite) sendTransactions(distribution TxTypeDistribution) {
	s.startTime = time.Now()

	for _, wallet := range s.wallets {
		var userStates []dto.UserStateWithID
		err := s.commander.Client().CallFor(&userStates, "hubble_getUserStates", []interface{}{wallet.PublicKey()})
		if err != nil {
			continue
		}

		fmt.Printf("%d states found for wallet %s\n", len(userStates), wallet.PublicKey().String())

		workers := 0
		for _, state := range userStates {
			s.stateIds = append(s.stateIds, state.StateID)

			s.waitGroup.Add(1)
			go s.runForWallet(wallet, state.StateID, distribution)

			workers += 1
			if workers >= maxConcurrentWorkers {
				break
			}
		}
	}

	s.Greater(len(s.stateIds), 0)

	s.waitGroup.Wait()
}

func (s *BenchmarkSuite) runForWallet(senderWallet bls.Wallet, senderStateID uint32, distribution TxTypeDistribution) {
	fmt.Printf("Starting worker on stateId %d address=%s\n", senderStateID, senderWallet.PublicKey().String())

	txsToWatch := make([]common.Hash, 0, maxQueuedBatchesCount)
	nonce := models.MakeUint256(0)

	for s.txsSent < txCount {
		// Send phase
		for len(txsToWatch) <= maxQueuedBatchesCount {
			var lastTxHash common.Hash

			for i := 0; i < txBatchSize; i++ {
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

				nonce = *nonce.AddN(1)
			}
			txsToWatch = append(txsToWatch, lastTxHash)
			atomic.AddInt64(&s.txsQueued, txBatchSize)
		}

		// Check phase
		newTxsToWatch := make([]common.Hash, 0)
		continueChecking := true
		for _, tx := range txsToWatch {
			if continueChecking {
				var receipt struct {
					Status txstatus.TransactionStatus
				}
				err := s.commander.Client().CallFor(&receipt, "hubble_getTransaction", []interface{}{tx})
				s.NoError(err)

				if receipt.Status != txstatus.Pending {
					atomic.AddInt64(&s.txsSent, txBatchSize)
					atomic.AddInt64(&s.txsQueued, -txBatchSize)
				} else {
					continueChecking = false
					newTxsToWatch = append(newTxsToWatch, tx)
				}
			} else {
				newTxsToWatch = append(newTxsToWatch, tx)
			}

			// If we send too many requests at the same time we can run out of OS ports
			time.Sleep(500 * time.Microsecond)
		}
		txsToWatch = newTxsToWatch

		// Report phase
		if s.lastReportedTxCount != s.txsSent {
			s.lastReportedTxCount = s.txsSent
			fmt.Printf(
				"Transfers sent: %d, throughput: %f tx/s, txs in queue: %d\n",
				s.txsSent,
				float64(s.txsSent)/(time.Since(s.startTime).Seconds()),
				s.txsQueued,
			)
		}
	}

	s.waitGroup.Done()
}

func (s *BenchmarkSuite) sendTransfer(wallet bls.Wallet, from, to uint32, nonce models.Uint256) common.Hash {
	transfer, err := api.SignTransfer(&wallet, dto.Transfer{
		FromStateID: &from,
		ToStateID:   &to,
		Amount:      models.NewUint256(1),
		Fee:         models.NewUint256(1),
		Nonce:       &nonce,
	})
	s.NoError(err)

	var transferHash common.Hash
	err = s.commander.Client().CallFor(&transferHash, "hubble_sendTransaction", []interface{}{*transfer})
	s.NoError(err)
	s.NotNil(transferHash)

	return transferHash
}

func (s *BenchmarkSuite) sendC2T(wallet bls.Wallet, from uint32, to *models.PublicKey, nonce models.Uint256) common.Hash {
	transfer, err := api.SignCreate2Transfer(&wallet, dto.Create2Transfer{
		FromStateID: &from,
		ToPublicKey: to,
		Amount:      models.NewUint256(1),
		Fee:         models.NewUint256(1),
		Nonce:       &nonce,
	})
	s.NoError(err)

	var transferHash common.Hash
	err = s.commander.Client().CallFor(&transferHash, "hubble_sendTransaction", []interface{}{*transfer})
	s.NoError(err)
	s.NotNil(transferHash)

	return transferHash
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

func randomInt(n int) int {
	//nolint:gosec
	return rand.Intn(n)
}

func randomFloat32() float32 {
	//nolint:gosec
	return rand.Float32()
}

func TestBenchmarkSuite(t *testing.T) {
	suite.Run(t, new(BenchmarkSuite))
}
