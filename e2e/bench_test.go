// +build e2e

package e2e

import (
	"fmt"
	"math/rand"
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
	nonces := map[uint32]*models.Uint256{}
	walletForState := map[uint32]bls.Wallet{}
	txsToWatch := map[uint32][]common.Hash{}

	for i := 0; i < 6; i++ {
		var userStates []dto.UserState
		err = commander.Client().CallFor(&userStates, "hubble_getUserStates", []interface{}{wallets[i].PublicKey()})
		require.NoError(t, err)
		for _, state := range userStates {
			stateIds = append(stateIds, state.StateID)
			nonces[state.StateID] = models.NewUint256(0)
			walletForState[state.StateID] = wallets[i]
			txsToWatch[state.StateID] = make([]common.Hash, 0)
		}
	}

	transfersSent := 0
	startTime := time.Now()

	// Total number of transactions to be sent.
	const txCount = 10_000

	// Number of transaction that will be sent in a single batch (unrelated to rollup "batches").
	const txBatchSize = 32

	// Maximum number of tx batches in queue.
	const maxQueuedBatchCount = 20

	for transfersSent < txCount {
		txSent := true
		for txSent {
			txSent = false

			for _, stateId := range stateIds {
				if len(txsToWatch[stateId]) > maxQueuedBatchCount { // max txs in queue
					continue
				}

				var hash *common.Hash
				for i := 0; i < txBatchSize; i++ {
					wallet := walletForState[stateId]
					nonce := nonces[stateId]
					to := stateId

					// Pick random receiver id thats different from sender's.
					for to == stateId {
						to = stateIds[rand.Intn(len(stateIds))]
					}

					hash = sendTransfer(t, commander, wallet, stateId, to, *nonce)
					if hash != nil {
						nonces[stateId] = nonces[stateId].AddN(1)
						txSent = true
					}
				}
				if hash != nil {
					txsToWatch[stateId] = append(txsToWatch[stateId], *hash)
				}
			}
		}

		txInQueue := 0
		for _, stateId := range stateIds {
			newTxsToWatch := make([]common.Hash, 0)
			skip := false
			for _, tx := range txsToWatch[stateId] {
				var sentTransfer dto.TransferReceipt
				if !skip {
					err = commander.Client().CallFor(&sentTransfer, "hubble_getTransaction", []interface{}{tx})
					require.NoError(t, err)
				}
				if skip || sentTransfer.Status == txstatus.Pending {
					newTxsToWatch = append(newTxsToWatch, tx)
					txInQueue += txBatchSize
					skip = true
				} else {
					transfersSent += txBatchSize
				}
			}
			txsToWatch[stateId] = newTxsToWatch
		}

		fmt.Printf("Transfers sent: %d, throughput: %f tx/s, txs in queue: %d\n", transfersSent, float64(transfersSent)/(time.Since(startTime).Seconds()), txInQueue)
	}
}

func sendTransfer(t *testing.T, commander Commander, wallet bls.Wallet, from uint32, to uint32, nonce models.Uint256) *common.Hash {
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

	return &transferHash
}
