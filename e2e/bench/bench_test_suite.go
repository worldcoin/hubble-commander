package bench

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BenchmarkConfig struct {
	// Total number of transactions to be sent.
	TxAmount int64

	// Number of transaction that will be sent in a single batch (unrelated to rollup "batches").
	TxBatchSize int64

	// Maximum number of tx batches in queue.
	MaxQueuedBatchesAmount int64

	// Maximum number of workers that send transactions.
	MaxConcurrentWorkers int64
}

type benchmarkTestSuite struct {
	*require.Assertions
	suite.Suite

	benchConfig BenchmarkConfig

	commander setup.Commander
	domain    bls.Domain
	wallets   []bls.Wallet
	stateIds  []uint32

	startTime time.Time
	waitGroup sync.WaitGroup

	// Only use atomic operations to increment those two counters.
	txsSent   int64
	txsQueued int64

	lastReportedTxAmount int64
}

func (s *benchmarkTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *benchmarkTestSuite) SetupTest(benchmarkConfig BenchmarkConfig) {
	s.SetupTestWithRollupConfig(benchmarkConfig, nil)
}

func (s *benchmarkTestSuite) SetupTestWithRollupConfig(benchmarkConfig BenchmarkConfig, cfg *config.RollupConfig) {
	s.benchConfig = benchmarkConfig

	commander, err := setup.NewConfiguredCommanderFromEnv(cfg)
	s.NoError(err)

	err = commander.Start()
	s.NoError(err)

	s.domain = e2e.GetDomain(s.T(), commander.Client())

	wallets, err := setup.CreateWallets(s.domain)
	s.NoError(err)

	s.commander = commander
	s.wallets = wallets
	s.stateIds = make([]uint32, 0)
	s.waitGroup = sync.WaitGroup{}
	s.txsSent = 0
	s.txsQueued = 0
	s.lastReportedTxAmount = 0
}

func (s *benchmarkTestSuite) TearDownTest() {
	s.NoError(s.commander.Stop())
}

func (s *benchmarkTestSuite) sendTransfer(wallet bls.Wallet, from, to uint32, nonce models.Uint256) common.Hash {
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

func (s *benchmarkTestSuite) sendC2T(wallet bls.Wallet, from uint32, to *models.PublicKey, nonce models.Uint256) common.Hash {
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

func (s *benchmarkTestSuite) sendMassMigration(wallet bls.Wallet, from uint32, nonce models.Uint256) common.Hash {
	massMigration, err := api.SignMassMigration(&wallet, dto.MassMigration{
		FromStateID: &from,
		SpokeID:     ref.Uint32(1),
		Amount:      models.NewUint256(1),
		Fee:         models.NewUint256(1),
		Nonce:       &nonce,
	})
	s.NoError(err)

	var massMigrationHash common.Hash
	err = s.commander.Client().CallFor(&massMigrationHash, "hubble_sendTransaction", []interface{}{*massMigration})
	s.NoError(err)
	s.NotNil(massMigrationHash)

	return massMigrationHash
}

func (s *benchmarkTestSuite) sendTransactions(
	walletAction func(senderWallet bls.Wallet, senderStateID uint32, nonce models.Uint256) common.Hash,
) {
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
			if state.Balance.CmpN(setup.InitialGenesisBalance) != 0 {
				continue
			}

			s.stateIds = append(s.stateIds, state.StateID)

			s.waitGroup.Add(1)
			go s.runForWallet(wallet, state.StateID, walletAction)

			workers += 1
			if workers >= int(s.benchConfig.MaxConcurrentWorkers) {
				break
			}
		}
	}

	s.Greater(len(s.stateIds), 0)

	s.waitGroup.Wait()
}

func (s *benchmarkTestSuite) runForWallet(
	senderWallet bls.Wallet,
	senderStateID uint32,
	action func(senderWallet bls.Wallet, senderStateID uint32, nonce models.Uint256) common.Hash,
) {
	fmt.Printf("Starting worker on stateId %d address=%s\n", senderStateID, senderWallet.PublicKey().String())

	txsToWatch := make([]common.Hash, 0, s.benchConfig.MaxQueuedBatchesAmount)
	nonce := models.MakeUint256(0)

	for s.txsSent < s.benchConfig.TxAmount {
		// Send phase
		for int64(len(txsToWatch)) <= s.benchConfig.MaxQueuedBatchesAmount {
			var lastTxHash common.Hash

			for i := 0; i < int(s.benchConfig.TxBatchSize); i++ {
				lastTxHash = action(senderWallet, senderStateID, nonce)
				nonce = *nonce.AddN(1)
			}
			txsToWatch = append(txsToWatch, lastTxHash)
			atomic.AddInt64(&s.txsQueued, s.benchConfig.TxBatchSize)
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
					atomic.AddInt64(&s.txsSent, s.benchConfig.TxBatchSize)
					atomic.AddInt64(&s.txsQueued, -s.benchConfig.TxBatchSize)
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
		if s.lastReportedTxAmount != s.txsSent {
			s.lastReportedTxAmount = s.txsSent
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
