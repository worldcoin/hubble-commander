//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	admintypes "github.com/Worldcoin/hubble-commander/client"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/healthstatus"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

// nolint:funlen
func TestCommanderMigrationMode(t *testing.T) {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = 4
	cfg.Rollup.MaxTxsPerCommitment = 10
	cfg.Rollup.MinCommitmentsPerBatch = 1
	if cfg.Ethereum == nil {
		log.Panicf("migration mode test cannot be run on simulator")
	}

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5001"
	cfg.Metrics.Port = "2001"
	activeCommander, err := setup.DeployAndCreateInProcessCommander(cfg, nil)
	require.NoError(t, err)

	adminRPCClient := testCreateAdminRPCClient(cfg)

	gethRPCClient, err := rpc.Dial(cfg.Ethereum.RPCURL)
	require.NoError(t, err)

	err = activeCommander.Start()
	require.NoError(t, err)
	defer func() {
		gethRPCClient.Close()
		require.NoError(t, activeCommander.Stop())
		require.NoError(t, os.Remove(*cfg.Bootstrap.ChainSpecPath))
	}()

	domain := GetDomain(t, activeCommander.Client())

	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	testStopMining(t, gethRPCClient)

	testSendInvalidTx(t, adminRPCClient, 4, wallets)
	testSendValidTxs(t, adminRPCClient, 0, 4, wallets, 1)
	testSendInvalidTx(t, adminRPCClient, 5, wallets)

	testWaitForBatchStatus(t, adminRPCClient, 1, batchstatus.Submitted)

	testConfigureCommander(t, adminRPCClient, dto.ConfigureParams{
		CreateBatches:      ref.Bool(false),
		AcceptTransactions: ref.Bool(true),
	})

	testSendValidTxs(t, adminRPCClient, 4, 4, wallets, 1)
	testSendValidTxs(t, adminRPCClient, 0, 4, wallets, 2)

	testConfigureCommander(t, adminRPCClient, dto.ConfigureParams{
		CreateBatches:      ref.Bool(false),
		AcceptTransactions: ref.Bool(false),
	})

	testValidateCommanderNotAcceptingTxs(t, adminRPCClient, wallets)

	cfg.Bootstrap.Prune = true
	cfg.Bootstrap.Migrate = true
	bootstrapURL := fmt.Sprint("http://localhost:", cfg.API.Port)
	cfg.Bootstrap.BootstrapNodeURL = &bootstrapURL
	cfg.API.Port = "5002"
	cfg.Metrics.Port = "2002"
	cfg.Badger.Path += "_migration"
	migratedCommander, err := setup.CreateInProcessCommander(cfg, nil)
	require.NoError(t, err)

	err = migratedCommander.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, migratedCommander.Stop())
	}()

	testWaitForCommanderReadyStatus(t, migratedCommander.Client())

	migratedAdminRPCClient := testCreateAdminRPCClient(cfg)

	testValidateMigration(t, migratedAdminRPCClient)

	testStartMining(t, gethRPCClient)

	// Wait for pending batch migrated from Commander 1 and validate it
	batch1 := testWaitForBatchStatus(t, migratedCommander.Client(), 1, batchstatus.Mined)
	testValidateBatch(t, migratedCommander.Client(), batch1, 4)

	// Wait for new batch created by Migrated Commander from pending transactions and validate it
	batch2 := testWaitForBatchStatus(t, migratedCommander.Client(), 2, batchstatus.Mined)
	testValidateBatch(t, migratedCommander.Client(), batch2, 8)

	// Send some txs to Migrated Commander and validate that it is creating new batches
	testSendValidTxs(t, migratedAdminRPCClient, 0, 7, wallets, 0)
	batch3 := testWaitForBatchStatus(t, migratedCommander.Client(), 3, batchstatus.Mined)
	testValidateBatch(t, migratedCommander.Client(), batch3, 7)
}

func testStopMining(t *testing.T, client *rpc.Client) {
	err := client.Call(nil, "miner_stop")
	require.NoError(t, err)
	log.Printf("Stopining mining blocks")
}

func testStartMining(t *testing.T, client *rpc.Client) {
	err := client.Call(nil, "miner_start")
	require.NoError(t, err)
	log.Printf("Starting mining blocks")
}

func testCreateAdminRPCClient(cfg *config.Config) jsonrpc.RPCClient {
	url := fmt.Sprint("http://localhost:", cfg.API.Port)
	adminRPCClient := jsonrpc.NewClientWithOpts(url, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			consts.AuthKeyHeader: cfg.API.AuthenticationKey,
		},
	})

	return adminRPCClient
}

func testSendValidTxs(t *testing.T, client jsonrpc.RPCClient, startingNonce, txsCount uint64, wallets []bls.Wallet, fromStateID uint32) {
	for i := uint64(0); i < txsCount; i++ {
		transfer, err := api.SignTransfer(&wallets[fromStateID], dto.Transfer{
			FromStateID: ref.Uint32(fromStateID),
			ToStateID:   ref.Uint32(fromStateID + 1),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(startingNonce + i),
		})
		require.NoError(t, err)

		var txHash common.Hash
		err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
		require.NoError(t, err)
		require.NotZero(t, txHash)
	}
}

func testSendInvalidTx(t *testing.T, client jsonrpc.RPCClient, nonce uint64, wallets []bls.Wallet) {
	transfer, err := api.SignTransfer(&wallets[1], dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(999),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(nonce),
	})
	require.NoError(t, err)

	var txHash common.Hash
	err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
	require.NoError(t, err)
	require.NotZero(t, txHash)
}

func testWaitForBatchStatus(
	t *testing.T,
	client jsonrpc.RPCClient,
	batchID uint64,
	status batchstatus.BatchStatus,
) *dto.BatchWithRootAndCommitments {
	var batch dto.BatchWithRootAndCommitments

	require.Eventually(t, func() bool {
		var rpcError *jsonrpc.RPCError
		err := client.CallFor(&batch, "hubble_getBatchByID", []interface{}{models.MakeUint256(batchID)})
		if err != nil && errors.As(err, &rpcError) {
			if rpcError.Code == 30000 {
				return false
			}
		}
		require.NoError(t, err)
		return batch.Status == status
	}, 30*time.Second, testutils.TryInterval)

	return &batch
}

func testValidateBatch(t *testing.T, client jsonrpc.RPCClient, batch *dto.BatchWithRootAndCommitments, txCount int) {
	require.Len(t, batch.Commitments, 1)

	commitmentID := models.CommitmentID{
		BatchID:      batch.ID,
		IndexInBatch: 0,
	}

	var commitment dto.TxCommitment
	err := client.CallFor(&commitment, "hubble_getCommitment", []interface{}{commitmentID})
	require.NoError(t, err)
	require.Len(t, commitment.Transactions, txCount)
}

func testConfigureCommander(t *testing.T, client jsonrpc.RPCClient, params dto.ConfigureParams) {
	_, err := client.Call("admin_configure", []interface{}{params})
	require.NoError(t, err)
}

func testWaitForCommanderReadyStatus(t *testing.T, client jsonrpc.RPCClient) {
	require.Eventually(t, func() bool {
		var status string
		var rpcError *jsonrpc.RPCError
		err := client.CallFor(&status, "hubble_getStatus")
		if err != nil && errors.As(err, &rpcError) {
			if rpcError.Code == 30000 {
				return false
			}
		}
		require.NoError(t, err)
		log.Printf(status)
		return status == healthstatus.Ready
	}, 30*time.Second, testutils.TryInterval)
}

func testValidateCommanderNotAcceptingTxs(t *testing.T, client jsonrpc.RPCClient, wallets []bls.Wallet) {
	transfer, err := api.SignTransfer(&wallets[0], dto.Transfer{
		FromStateID: ref.Uint32(0),
		ToStateID:   ref.Uint32(1),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	})
	require.NoError(t, err)

	var txHash common.Hash
	var rpcError *jsonrpc.RPCError
	err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
	require.True(t, errors.As(err, &rpcError))
	require.Equal(t, 10017, rpcError.Code)
}

func testValidateMigration(t *testing.T, client jsonrpc.RPCClient) {
	testValidateFailedTxs(t, client)
	testValidatePendingTxs(t, client)
	testValidatePendingBatches(t, client)
}

func testValidatePendingTxs(t *testing.T, client jsonrpc.RPCClient) {
	pendingTxs := make([]admintypes.Transaction, 0)
	err := client.CallFor(&pendingTxs, "admin_getPendingTransactions")
	require.NoError(t, err)
	require.Len(t, pendingTxs, 8)
}

func testValidateFailedTxs(t *testing.T, client jsonrpc.RPCClient) {
	failedTxs := make([]admintypes.Transaction, 0)
	err := client.CallFor(&failedTxs, "admin_getFailedTransactions")
	require.NoError(t, err)
	require.Len(t, failedTxs, 2)
}

func testValidatePendingBatches(t *testing.T, client jsonrpc.RPCClient) {
	require.Eventually(t, func() bool {
		var pendingBatches []admintypes.Batch
		err := client.CallFor(&pendingBatches, "admin_getPendingBatches")
		if err != nil {
			log.Debugf(err.Error())
			return false
		}
		require.NoError(t, err)

		require.Len(t, pendingBatches, 1)
		require.Len(t, pendingBatches[0].Commitments, 1)
		require.Len(t, pendingBatches[0].Commitments[0].Transactions, 4)

		return true
	}, 30*time.Second, testutils.TryInterval)
}
