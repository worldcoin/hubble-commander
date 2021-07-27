// +build e2e

package e2e

import (
	"bytes"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestMeasureDisputeGasUsage(t *testing.T) {
	cmd, err := setup.NewCommanderFromEnv(true)
	require.NoError(t, err)
	err = cmd.Start()
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
	measureDisputeTransitionC2TInvalidStateRoot(t, cmd.Client(), ethClient, wallets)
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

func measureDisputeTransitionC2TInvalidStateRoot(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, wallets []bls.Wallet) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	send32C2TBatchWithInvalidStateRoot(t, ethClient, wallets)
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

	sendTransferCommitment(t, ethClient, bytes.Repeat(encodedTransfer, 32), 2)
}

func send32C2TBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client, wallets []bls.Wallet) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(38),
	}

	registrations, unsubscribe, err := ethClient.WatchBatchAccountRegistrations(&bind.WatchOpts{})
	require.NoError(t, err)
	defer unsubscribe()

	publicKeyBatch := make([]models.PublicKey, 16)
	registeredPubKeyIDs := make([]uint32, 0, 32)
	walletIndex := len(wallets) - 32
	for i := 0; i < 2; i++ {
		for j := range publicKeyBatch {
			publicKeyBatch[j] = *wallets[walletIndex].PublicKey()
			walletIndex++
		}
		pubKeyIDs, err := ethClient.RegisterBatchAccount(publicKeyBatch, registrations)
		require.NoError(t, err)
		registeredPubKeyIDs = append(registeredPubKeyIDs, pubKeyIDs...)
	}

	encodedTransfers := make([]byte, 0, encoder.Create2TransferLength*32)
	for i := range registeredPubKeyIDs {
		transfer.ToStateID = ref.Uint32(uint32(38 + i))
		encodedTx, err := encoder.EncodeCreate2TransferForCommitment(&transfer, registeredPubKeyIDs[i])
		require.NoError(t, err)

		encodedTransfers = append(encodedTransfers, encodedTx...)
	}

	sendC2TCommitment(t, ethClient, encodedTransfers, 2)
}
