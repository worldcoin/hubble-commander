// +build e2e

package e2e

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
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

type BenchmarkSuite struct {
	*require.Assertions
	suite.Suite

	commander setup.Commander
	wallets   []bls.Wallet
	stateIds  []uint32

	startTime time.Time
	waitGroup sync.WaitGroup

	// Only use atomic operations to increment those two counters.
	transfersSent int64
	txsQueued     int64
}

func (s *BenchmarkSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	commander, err := setup.NewCommanderFromEnv(true)
	s.NoError(err)

	err = commander.Start()
	s.NoError(err)

	domain := getDomain(s.T(), commander.Client())
	wallets, err := setup.CreateWallets(domain)
	s.NoError(err)

	s.commander = commander
	s.wallets = wallets
	s.stateIds = make([]uint32, 0)
}

func (s *BenchmarkSuite) TearDownSuite() {
	s.NoError(s.commander.Stop())
}

func (s *BenchmarkSuite) TestBenchCommander() {
	s.startTime = time.Now()

	for _, wallet := range s.wallets {
		var userStates []dto.UserState
		err := s.commander.Client().CallFor(&userStates, "hubble_getUserStates", []interface{}{wallet.PublicKey()})
		if err != nil {
			continue
		}

		fmt.Printf("%d states found for wallet %s\n", len(userStates), wallet.PublicKey().String())

		for _, state := range userStates {
			s.stateIds = append(s.stateIds, state.StateID)

			s.waitGroup.Add(1)
			go s.runForWallet(wallet, state.StateID)
		}
	}

	s.Greater(len(s.stateIds), 0)

	s.waitGroup.Wait()
}

func (s *BenchmarkSuite) runForWallet(senderWallet bls.Wallet, senderStateId uint32) {
	fmt.Printf("Starting worker on stateId %d address=%s\n", senderStateId, senderWallet.PublicKey().String())

	txsToWatch := make([]common.Hash, 0, maxQueuedBatchesCount)
	nonce := models.MakeUint256(0)

	for s.transfersSent < txCount {

		// Send phase
		for len(txsToWatch) <= maxQueuedBatchesCount {
			var lastTxHash common.Hash
			for i := 0; i < txBatchSize; i++ {

				// Pick random receiver id that's different from sender's.
				to := s.stateIds[rand.Intn(len(s.stateIds))]
				for to == senderStateId {
					to = s.stateIds[rand.Intn(len(s.stateIds))]
				}

				lastTxHash = s.sendTransfer(senderWallet, senderStateId, to, nonce)
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
				var sentTransfer dto.TransferReceipt
				err := s.commander.Client().CallFor(&sentTransfer, "hubble_getTransaction", []interface{}{tx})
				s.NoError(err)

				if sentTransfer.Status != txstatus.Pending {
					atomic.AddInt64(&s.transfersSent, txBatchSize)
					atomic.AddInt64(&s.txsQueued, -txBatchSize)
				} else {
					continueChecking = false
					newTxsToWatch = append(newTxsToWatch, tx)
				}
			} else {
				newTxsToWatch = append(newTxsToWatch, tx)
			}
		}
		txsToWatch = newTxsToWatch

		// Report phase
		fmt.Printf("Transfers sent: %d, throughput: %f tx/s, txs in queue: %d\n", s.transfersSent, float64(s.transfersSent)/(time.Since(s.startTime).Seconds()), s.txsQueued)

	}

	s.waitGroup.Done()
}

func (s *BenchmarkSuite) sendTransfer(wallet bls.Wallet, from uint32, to uint32, nonce models.Uint256) common.Hash {
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

func TestBenchmarkSuite(t *testing.T) {
	suite.Run(t, new(BenchmarkSuite))
}
