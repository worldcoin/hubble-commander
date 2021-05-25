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
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBenchCommander(t *testing.T) {
	commander, err := NewCommanderFromEnv(true)
	require.NoError(t, err)
	err = commander.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, commander.Stop())
	}()

	domain := getDomain(t, commander.Client())

	wallets, err := createWallets(domain)
	require.NoError(t, err)

	var version string
	err = commander.Client().CallFor(&version, "hubble_getVersion")
	require.NoError(t, err)
	require.Equal(t, config.GetConfig().API.Version, version)

	stateIds := make([]uint32, 0)

	// Total number of transactions to be sent.
	const txCount = 10_000

	// Number of transaction that will be sent in a single batch (unrelated to rollup "batches").
	const txBatchSize = 32

	// Maximum number of tx batches in queue.
	const maxQueuedBatchCount = 20

	// Only use atomic operations to increment those two counters.
	transfersSent := int64(0)
	txsQueued := int64(0)

	startTime := time.Now()

	var waitGroup sync.WaitGroup

	runForWallet := func(wallet bls.Wallet, stateId uint32) {
		fmt.Printf("Starting worker on stateId %d address=%s\n", stateId, wallet.PublicKey().String())

		txsToWatch := make([]common.Hash, 0, maxQueuedBatchCount)
		nonce := models.MakeUint256(0)

		for transfersSent < txCount {

			// Send phase
			for len(txsToWatch) <= maxQueuedBatchCount {
				var hash common.Hash
				for i := 0; i < txBatchSize; i++ {

					// Pick random receiver id thats different from sender's.
					to := stateIds[rand.Intn(len(stateIds))]
					for to == stateId {
						to = stateIds[rand.Intn(len(stateIds))]
					}

					hash = sendTransfer(t, commander, wallet, stateId, to, nonce)
					nonce = *nonce.AddN(1)
				}
				txsToWatch = append(txsToWatch, hash)
				atomic.AddInt64(&txsQueued, txBatchSize)
			}

			// Check phase
			newTxsToWatch := make([]common.Hash, 0)
			skip := false
			for _, tx := range txsToWatch {
				var sentTransfer dto.TransferReceipt
				if !skip {
					err = commander.Client().CallFor(&sentTransfer, "hubble_getTransaction", []interface{}{tx})
					require.NoError(t, err)
				}
				if skip || sentTransfer.Status == txstatus.Pending {
					newTxsToWatch = append(newTxsToWatch, tx)
					skip = true
				} else {
					atomic.AddInt64(&transfersSent, txBatchSize)
					atomic.AddInt64(&txsQueued, -txBatchSize)

				}
			}
			txsToWatch = newTxsToWatch

			// Report phase
			fmt.Printf("Transfers sent: %d, throughput: %f tx/s, txs in queue: %d\n", transfersSent, float64(transfersSent)/(time.Since(startTime).Seconds()), txsQueued)
		}

		waitGroup.Done()
	}

	for _, wallet := range wallets {
		var userStates []dto.UserState
		err = commander.Client().CallFor(&userStates, "hubble_getUserStates", []interface{}{wallet.PublicKey()})
		if err != nil {
			continue
		}

		fmt.Printf("%d states found for wallet %s\n", len(userStates), wallet.PublicKey().String())

		for _, state := range userStates {
			stateIds = append(stateIds, state.StateID)

			waitGroup.Add(1)
			go runForWallet(wallet, state.StateID)
		}
	}

	require.Greater(t, len(stateIds), 0)

	waitGroup.Wait()
}

func sendTransfer(t *testing.T, commander Commander, wallet bls.Wallet, from uint32, to uint32, nonce models.Uint256) common.Hash {
	transfer, err := api.SignTransfer(&wallet, dto.Transfer{
		FromStateID: &from,
		ToStateID:   &to,
		Amount:      models.NewUint256(1),
		Fee:         models.NewUint256(1),
		Nonce:       &nonce,
	})
	require.NoError(t, err)

	var transferHash common.Hash
	err = commander.Client().CallFor(&transferHash, "hubble_sendTransaction", []interface{}{*transfer})
	require.NoError(t, err)
	require.NotNil(t, transferHash)

	return transferHash
}
