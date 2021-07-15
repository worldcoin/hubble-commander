package e2e

import (
	"testing"
	"time"

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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/event"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommanderDispute(t *testing.T) {
	cmd := setup.CreateInProcessCommander()
	err := cmd.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cmd.Stop())
	}()

	domain := getDomain(t, cmd.Client())

	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	senderWallet := wallets[1]

	firstTransferHash := testSendTransfer(t, cmd.Client(), senderWallet, models.NewUint256(0))
	testGetTransaction(t, cmd.Client(), firstTransferHash)
	send31MoreTransfers(t, cmd.Client(), senderWallet)

	waitForTxToBeIncludedInBatch(t, cmd.Client(), firstTransferHash)

	ethClient := newEthClient(t, cmd.Client())

	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	sendInvalidBatch(t, ethClient)
	waitForRollbackToFinish(t, sink, subscription)

	testBatchesAfterDispute(t, cmd.Client())
}

func waitForRollbackToFinish(
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

func sendInvalidBatch(t *testing.T, ethClient *eth.Client) {
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

	_, err = ethClient.GetBatch(models.NewUint256(2))
	require.NoError(t, err)
}

func testBatchesAfterDispute(t *testing.T, client jsonrpc.RPCClient) {
	var batches []dto.Batch
	err := client.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})

	require.NoError(t, err)
	require.Len(t, batches, 1)
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
