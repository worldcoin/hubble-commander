package bench

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

type BenchmarkConfig struct {
	// Total number of transactions to be sent.
	TxCount int64

	// Number of transaction that will be sent in a single batch (unrelated to rollup "batches").
	TxBatchSize int64

	// Maximum number of tx batches in queue.
	MaxQueuedBatchesCount int64

	// Maximum number of workers that send transactions.
	MaxConcurrentWorkers int64
}

type benchmarkTestSuite struct {
	setup.E2ETestSuite

	benchConfig BenchmarkConfig

	stateIds []uint32

	startTime time.Time
	waitGroup sync.WaitGroup

	// Only use atomic operations to increment those two counters.
	txsSent   int64
	txsQueued int64

	lastReportedTxCount int64
}

func (s *benchmarkTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *benchmarkTestSuite) SetupTest(benchmarkConfig BenchmarkConfig) {
	s.SetupTestEnvironment(nil, nil)

	s.benchConfig = benchmarkConfig
	s.stateIds = make([]uint32, 0)
	s.waitGroup = sync.WaitGroup{}
	s.txsSent = 0
	s.txsQueued = 0
	s.lastReportedTxCount = 0
}

func (s *benchmarkTestSuite) sendTransferFromCustomWallet(wallet bls.Wallet, from, to uint32, nonce models.Uint256) common.Hash {
	transfer, err := api.SignTransfer(&wallet, dto.Transfer{
		FromStateID: &from,
		ToStateID:   &to,
		Amount:      models.NewUint256(1),
		Fee:         models.NewUint256(1),
		Nonce:       &nonce,
	})
	s.NoError(err)

	var transferHash common.Hash
	err = s.RPCClient.CallFor(&transferHash, "hubble_sendTransaction", []interface{}{*transfer})
	s.NoError(err)
	s.NotNil(transferHash)

	return transferHash
}

func (s *benchmarkTestSuite) sendC2TFromCustomWallet(
	wallet bls.Wallet,
	from uint32,
	to *models.PublicKey,
	nonce models.Uint256,
) common.Hash {
	transfer, err := api.SignCreate2Transfer(&wallet, dto.Create2Transfer{
		FromStateID: &from,
		ToPublicKey: to,
		Amount:      models.NewUint256(1),
		Fee:         models.NewUint256(1),
		Nonce:       &nonce,
	})
	s.NoError(err)

	var transferHash common.Hash
	err = s.RPCClient.CallFor(&transferHash, "hubble_sendTransaction", []interface{}{*transfer})
	s.NoError(err)
	s.NotNil(transferHash)

	return transferHash
}

func (s *benchmarkTestSuite) sendMMFromCustomWallet(wallet bls.Wallet, from uint32, nonce models.Uint256) common.Hash {
	massMigration, err := api.SignMassMigration(&wallet, dto.MassMigration{
		FromStateID: &from,
		SpokeID:     ref.Uint32(1),
		Amount:      models.NewUint256(1),
		Fee:         models.NewUint256(1),
		Nonce:       &nonce,
	})
	s.NoError(err)

	var massMigrationHash common.Hash
	err = s.RPCClient.CallFor(&massMigrationHash, "hubble_sendTransaction", []interface{}{*massMigration})
	s.NoError(err)
	s.NotNil(massMigrationHash)

	return massMigrationHash
}

func (s *benchmarkTestSuite) sendTransactions(
	walletAction func(senderWallet bls.Wallet, senderStateID uint32, nonce models.Uint256) common.Hash,
) {
	s.startTime = time.Now()

	for _, wallet := range s.Wallets {
		var userStates []dto.UserStateWithID
		err := s.RPCClient.CallFor(&userStates, "hubble_getUserStates", []interface{}{wallet.PublicKey()})
		if err != nil {
			continue
		}

		fmt.Printf("%d states found for wallet %s\n", len(userStates), wallet.PublicKey().String())

		workers := 0
		for _, state := range userStates {
			if !state.Balance.EqN(setup.InitialGenesisBalance) {
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

	txsToWatch := make([]common.Hash, 0, s.benchConfig.MaxQueuedBatchesCount)
	nonce := models.MakeUint256(0)

	for s.txsSent < s.benchConfig.TxCount {
		// Send phase
		for int64(len(txsToWatch)) <= s.benchConfig.MaxQueuedBatchesCount {
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
				receipt := s.GetTransaction(tx)

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
