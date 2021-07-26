package e2e

import (
	"bytes"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestMeasureDisputeGasUsage(t *testing.T) {
	//cmd, err := setup.NewCommanderFromEnv(true)
	//require.NoError(t, err)
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

	ethClient := newEthClient(t, cmd.Client())

	testSendTransferBatch(t, cmd.Client(), senderWallet, 0)

	measureDisputeTransitionTransferInvalidStateRoot(t, cmd.Client(), ethClient)
	//testDisputeTransitionC2TInvalidStateRoot(t, cmd.Client(), ethClient, wallets[len(wallets)-1])
}

func measureDisputeTransitionTransferInvalidStateRoot(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	send32TransfersBatchWithInvalidStateRoot(t, ethClient)
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 1)
}

func send32TransfersBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client) {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	require.NoError(t, err)

	commitment := models.Commitment{
		Transactions:      bytes.Repeat(encodedTransfer, 32),
		FeeReceiver:       0,
		CombinedSignature: models.Signature{},
		PostStateRoot:     utils.RandomHash(),
	}
	commitments := make([]models.Commitment, 32)
	for i := range commitments {
		commitments[i] = commitment
	}
	submitTransfersBatch(t, ethClient, commitments, 2)
}

func submitTransfersBatch(t *testing.T, ethClient *eth.Client, commitments []models.Commitment, batchID uint64) {
	transaction, err := ethClient.SubmitTransfersBatch(commitments)
	require.NoError(t, err)

	waitForSubmittedBatch(t, ethClient, transaction, batchID)
}
