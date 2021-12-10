//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ybbus/jsonrpc/v2"
)

func testSubmitTransferBatch(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, startNonce uint64) common.Hash {
	firstTransferHash := testSendTransfer(t, client, senderWallet, startNonce)
	testGetTransaction(t, client, firstTransferHash)
	send31MoreTransfers(t, client, senderWallet, startNonce+1)
	return firstTransferHash
}

func testSubmitC2TBatch(
	t *testing.T,
	client jsonrpc.RPCClient,
	senderWallet bls.Wallet,
	wallets []bls.Wallet,
	targetPublicKey *models.PublicKey,
	startNonce uint64,
) common.Hash {
	firstTransferHash := testSendCreate2Transfer(t, client, senderWallet, targetPublicKey, startNonce)
	testGetTransaction(t, client, firstTransferHash)
	send31MoreCreate2Transfers(t, client, senderWallet, wallets, startNonce+1)
	return firstTransferHash
}

func testSubmitMassMigrationBatch(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, startNonce uint64) common.Hash {
	firstMassMigrationHash := testSendMassMigration(t, client, senderWallet, startNonce)
	testGetTransaction(t, client, firstMassMigrationHash)
	send31MoreMassMigrations(t, client, senderWallet, startNonce+1)
	return firstMassMigrationHash
}

func submitTxBatchAndWait(t *testing.T, client jsonrpc.RPCClient, submit func() common.Hash) {
	firstTxHash := submit()
	waitForTxToBeIncludedInBatch(t, client, firstTxHash)
}
