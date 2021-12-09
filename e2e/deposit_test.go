//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/erc20"
	"github.com/Worldcoin/hubble-commander/contracts/test/customtoken"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

const queueDepositGasLimit = 600_000

func testSubmitDepositBatchAndWait(t *testing.T, client jsonrpc.RPCClient) {
	makeDeposits(t, client)
	waitForBatch(t, client, models.MakeUint256(4))
}

func makeDeposits(t *testing.T, client jsonrpc.RPCClient) {
	ethClient := newEthClient(t, client)

	tokenAddress := deployExampleToken(t, ethClient)
	tokenID := registerToken(t, ethClient, tokenAddress)
	approveToken(t, ethClient, tokenAddress)
	amount := models.NewUint256FromBig(*utils.ParseEther("10"))

	subtreeDepth, err := ethClient.GetMaxSubTreeDepthParam()
	require.NoError(t, err)
	depositCount := 1 << *subtreeDepth
	txs := make([]types.Transaction, 0, depositCount)
	for i := 0; i < depositCount; i++ {
		var tx *types.Transaction
		tx, err = ethClient.QueueDeposit(queueDepositGasLimit, models.NewUint256(1), amount, tokenID)
		require.NoError(t, err)
		txs = append(txs, *tx)
	}
	_, err = chain.WaitForMultipleTxs(ethClient.Blockchain.GetBackend(), txs...)
	require.NoError(t, err)
}

func registerToken(t *testing.T, ethClient *eth.Client, tokenAddress common.Address) *models.Uint256 {
	err := ethClient.RequestRegisterTokenAndWait(tokenAddress)
	require.NoError(t, err)

	tokenID, err := ethClient.FinalizeRegisterTokenAndWait(tokenAddress)
	require.NoError(t, err)

	return tokenID
}

func approveToken(t *testing.T, ethClient *eth.Client, tokenAddress common.Address) {
	token, err := erc20.NewERC20(tokenAddress, ethClient.Blockchain.GetBackend())
	require.NoError(t, err)

	tx, err := token.Approve(ethClient.Blockchain.GetAccount(), ethClient.ChainState.DepositManager, utils.ParseEther("100"))
	require.NoError(t, err)

	_, err = chain.WaitToBeMined(ethClient.Blockchain.GetBackend(), tx)
	require.NoError(t, err)
}

func deployExampleToken(t *testing.T, ethClient *eth.Client) common.Address {
	tokenAddress, tx, _, err := customtoken.DeployTestCustomToken(
		ethClient.Blockchain.GetAccount(),
		ethClient.Blockchain.GetBackend(),
		"ExampleToken",
		"EXP",
	)
	require.NoError(t, err)

	_, err = chain.WaitToBeMined(ethClient.Blockchain.GetBackend(), tx)
	require.NoError(t, err)

	return tokenAddress
}

func waitForBatch(t *testing.T, client jsonrpc.RPCClient, batchID models.Uint256) {
	require.Eventually(t, func() bool {
		var batch dto.BatchWithRootAndCommitments
		var rpcError *jsonrpc.RPCError
		err := client.CallFor(&batch, "hubble_getBatchByID", []interface{}{batchID})
		if err != nil && errors.As(err, &rpcError) {
			if rpcError.Code == 30000 {
				return false
			}
		}
		require.NoError(t, err)
		return true
	}, 30*time.Second, testutils.TryInterval)
}
