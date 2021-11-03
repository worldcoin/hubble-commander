//go:build e2e
// +build e2e

package bench

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/ethereum/go-ethereum/common"
)

func (s *BenchmarkSuite) registerPubKeys() {
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
			go s.runRegistrationsForWallet(wallet, state.StateID)

			workers += 1
			if workers >= maxConcurrentWorkers {
				break
			}
		}
	}

	s.Greater(len(s.stateIds), 0)

	s.waitGroup.Wait()
}

func (s *BenchmarkSuite) runRegistrationsForWallet(senderWallet bls.Wallet, senderStateID uint32) {
	fmt.Printf("Starting worker on stateId %d address=%s\n", senderStateID, senderWallet.PublicKey().String())

	txsToWatch := make([]common.Hash, 0, maxQueuedBatchesCount)
	nonce := models.MakeUint256(0)

	for s.txsSent < txCount {
		// Send phase
		for len(txsToWatch) <= maxQueuedBatchesCount {
			var lastTxHash common.Hash

			for i := 0; i < txBatchSize; i++ {
				to := s.unregisteredWallets[s.unregisteredWalletsIndex].PublicKey()
				s.unregisteredWalletsIndex++

				lastTxHash = s.sendC2T(senderWallet, senderStateID, to, nonce)

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
