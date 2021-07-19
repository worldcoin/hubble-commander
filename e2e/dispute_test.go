package e2e

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/event"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommanderDispute(t *testing.T) {
	cmd := setup.CreateInProcessCommander()
	//cmd, err := setup.NewCommanderFromEnv(true)
	//require.NoError(t, err)
	err := cmd.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cmd.Stop())
	}()

	domain := getDomain(t, cmd.Client())
	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	senderWallet := wallets[1]

	ethClient := newEthClient(t, cmd.Client())

	testDisputeTransitionTransfer(t, cmd.Client(), ethClient, senderWallet)
}

func testDisputeTransitionTransfer(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, senderWallet bls.Wallet) {
	testSendBatch(t, client, senderWallet, 0)

	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	sendInvalidTransferBatch(t, ethClient)
	testRollbackCompletion(t, sink, subscription)

	testBatchesAfterDispute(t, client)

	testSendBatch(t, client, senderWallet, 32)
}

func testSendBatch(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, startNonce uint64) {
	firstTransferHash := testSendTransfer(t, client, senderWallet, startNonce)
	testGetTransaction(t, client, firstTransferHash)
	send31MoreTransfers(t, client, senderWallet, startNonce+1)

	waitForTxToBeIncludedInBatch(t, client, firstTransferHash)
}

func testSendC2TBatch(
	t *testing.T,
	client jsonrpc.RPCClient,
	senderWallet bls.Wallet,
	wallets []bls.Wallet,
	targetPublicKey *models.PublicKey,
	startNonce uint64,
) {
	firstTransferHash := testSendCreate2Transfer(t, client, senderWallet, targetPublicKey, startNonce)
	testGetTransaction(t, client, firstTransferHash)
	send31MoreCreate2Transfers(t, client, senderWallet, wallets, startNonce+1)

	waitForTxToBeIncludedInBatch(t, client, firstTransferHash)
}

func testBatchesAfterDispute(t *testing.T, client jsonrpc.RPCClient) {
	var batches []dto.Batch
	err := client.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})

	require.NoError(t, err)
	require.Len(t, batches, 1)
}

func testRollbackCompletion(
	t *testing.T,
	sink chan *rollup.RollupRollbackStatus,
	subscription event.Subscription,
) {
	require.Eventually(t, func() bool {
		for {
			select {
			case err := <-subscription.Err():
				require.NoError(t, err)
				return false
			case rollbackStatus := <-sink:
				if rollbackStatus.Completed {
					return true
				}
			}
		}
	}, 30*time.Second, testutils.TryInterval)
}

func sendInvalidTransferBatch(t *testing.T, ethClient *eth.Client) {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(2_000_000_000_000_000_000),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	require.NoError(t, err)

	sendCommitment(t, ethClient, encodedTransfer, 2)
}

func sendInvalidCreate2TransferBatch(t *testing.T, ethClient *eth.Client, toPublicKey *models.PublicKey) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(2_000_000_000_000_000_000),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(6),
	}

	registrations, unsubscribe, err := ethClient.WatchRegistrations(&bind.WatchOpts{})
	require.NoError(t, err)
	defer unsubscribe()

	pubKeyID, err := ethClient.RegisterAccount(toPublicKey, registrations)
	require.NoError(t, err)

	encodedTransfer, err := encoder.EncodeCreate2TransferForCommitment(&transfer, *pubKeyID)
	require.NoError(t, err)

	sendCommitment(t, ethClient, encodedTransfer, 3)
}

func sendCommitment(t *testing.T, ethClient *eth.Client, encodedTransfer []byte, batchID uint64) {
	commitment := models.Commitment{
		Transactions:      encodedTransfer,
		FeeReceiver:       0,
		CombinedSignature: models.Signature{},
		PostStateRoot:     utils.RandomHash(),
	}
	transaction, err := ethClient.SubmitTransfersBatch([]models.Commitment{commitment})
	require.NoError(t, err)

	_, err = deployer.WaitToBeMined(ethClient.ChainConnection.GetBackend(), transaction)
	require.NoError(t, err)

	_, err = ethClient.GetBatch(models.NewUint256(batchID))
	require.NoError(t, err)
}

func newEthClient(t *testing.T, client jsonrpc.RPCClient) *eth.Client {
	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	require.NoError(t, err)

	chainState := models.ChainState{
		ChainID:         info.ChainID,
		AccountRegistry: info.AccountRegistry,
		DeploymentBlock: info.DeploymentBlock,
		Rollup:          info.Rollup,
	}

	cfg := config.GetConfig()
	chain, err := deployer.NewRPCChainConnection(cfg.Ethereum)
	require.NoError(t, err)

	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, chain.GetBackend())
	require.NoError(t, err)

	rollupContract, err := rollup.NewRollup(chainState.Rollup, chain.GetBackend())
	require.NoError(t, err)

	ethClient, err := eth.NewClient(chain, &eth.NewClientParams{
		ChainState:      chainState,
		Rollup:          rollupContract,
		AccountRegistry: accountRegistry,
	})
	require.NoError(t, err)
	return ethClient
}
