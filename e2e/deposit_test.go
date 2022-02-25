//go:build e2e
// +build e2e

package e2e

import (
	testing "testing"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/test/customtoken"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

const queueDepositGasLimit = 600_000

func testSubmitDepositBatchAndWait(
	t *testing.T,
	cmd setup.Commander,
	ethClient *eth.Client,
	token *models.RegisteredToken,
	batchID uint64,
) {
	// wait for previous batch to be mined
	waitForBatch(t, cmd.Client(), models.MakeUint256(batchID-1))
	makeDeposits(t, ethClient, token)
	waitForBatch(t, cmd.Client(), models.MakeUint256(batchID))
}

func makeDeposits(t *testing.T, ethClient *eth.Client, token *models.RegisteredToken) {
	amount := models.NewUint256FromBig(*utils.ParseEther("10"))

	subtreeDepth, err := ethClient.GetMaxSubtreeDepthParam()
	require.NoError(t, err)
	depositCount := 1 << *subtreeDepth
	txs := make([]types.Transaction, 0, depositCount)
	for i := 0; i < depositCount; i++ {
		var tx *types.Transaction
		tx, err = ethClient.QueueDeposit(queueDepositGasLimit, models.NewUint256(1), amount, &token.ID)
		require.NoError(t, err)
		txs = append(txs, *tx)
	}
	_, err = ethClient.WaitForMultipleTxs(txs...)
	require.NoError(t, err)
}

func approveTokens(t *testing.T, token *customtoken.TestCustomToken, ethClient *eth.Client) {
	tx, err := token.Approve(ethClient.Blockchain.GetAccount(), ethClient.ChainState.DepositManager, utils.ParseEther("100"))
	require.NoError(t, err)

	_, err = ethClient.WaitToBeMined(tx)
	require.NoError(t, err)
}

func getDeployedToken(t *testing.T, ethClient *eth.Client) (*models.RegisteredToken, *customtoken.TestCustomToken) {
	registeredToken, err := ethClient.GetRegisteredToken(models.NewUint256(0))
	require.NoError(t, err)

	tokenContract, err := customtoken.NewTestCustomToken(registeredToken.Contract, ethClient.Blockchain.GetBackend())
	require.NoError(t, err)

	return registeredToken, tokenContract
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
		return batch.Status != batchstatus.Submitted
	}, 30*time.Second, testutils.TryInterval)
}
