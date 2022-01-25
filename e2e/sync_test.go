//go:build e2e
// +build e2e

package e2e

import (
	"os"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/testutils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCommanderSync(t *testing.T) {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = 32
	cfg.Rollup.MaxTxsPerCommitment = 32
	cfg.Rollup.MinCommitmentsPerBatch = 1
	if cfg.Ethereum.RPCURL == "simulator" {
		// newEthClient attempts to connect to a geth node which does not exist
		log.Panicf("sync test cannot be run on simulator")
	}

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5001"
	cfg.Metrics.Port = "2001"
	activeCommander, err := setup.DeployAndCreateInProcessCommander(cfg, nil)
	require.NoError(t, err)

	err = activeCommander.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, activeCommander.Stop())
		require.NoError(t, os.Remove(*cfg.Bootstrap.ChainSpecPath))
	}()

	domain := GetDomain(t, activeCommander.Client())

	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	feeReceiverWallet := wallets[0]
	senderWallet := wallets[1]

	testGetVersion(t, activeCommander.Client())
	testGetUserStates(t, activeCommander.Client(), senderWallet)

	submitBatchesAndWait(t, activeCommander, senderWallet, wallets)

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5002"
	cfg.Metrics.Port = "2002"
	cfg.Badger.Path += "_passive"
	cfg.Ethereum.PrivateKey = "ab6919fd6ac00246bb78657e0696cf72058a4cb395133d074eabaddb83d8b00c"
	passiveCommander, err := setup.CreateInProcessCommander(cfg, nil)
	require.NoError(t, err)

	err = passiveCommander.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, passiveCommander.Stop())
	}()

	var networkInfo dto.NetworkInfo
	err = activeCommander.Client().CallFor(&networkInfo, "hubble_getNetworkInfo")
	require.NoError(t, err)

	latestBatch := networkInfo.LatestBatch

	require.Eventually(t, func() bool {
		var networkInfo dto.NetworkInfo
		err := passiveCommander.Client().CallFor(&networkInfo, "hubble_getNetworkInfo")
		require.NoError(t, err)
		return networkInfo.LatestBatch != nil && networkInfo.LatestBatch.Cmp(latestBatch) >= 0
	}, 30*time.Second, testutils.TryInterval)

	testSenderStateAfterTransfers(
		t,
		passiveCommander.Client(),
		senderWallet,
		32+32+32,
		32*100+32*100+32*100,
	)

	testFeeReceiverStateAfterTransfers(
		t,
		passiveCommander.Client(),
		feeReceiverWallet,
		32*10*3,
	)

	testGetBatches(t, passiveCommander.Client(), 5)
}

func submitBatchesAndWait(t *testing.T, activeCommander *setup.InProcessCommander, senderWallet bls.Wallet, wallets []bls.Wallet) {
	firstTransferHash := testSubmitTransferBatch(t, activeCommander.Client(), senderWallet, 0)
	firstC2THash := testSubmitC2TBatch(t, activeCommander.Client(), senderWallet, wallets, wallets[len(wallets)-32].PublicKey(), 32)
	firstMMHash := testSubmitMassMigrationBatch(t, activeCommander.Client(), senderWallet, 64)
	testSubmitDepositBatchAndWait(t, activeCommander.Client(), 4)

	waitForTxToBeIncludedInBatch(t, activeCommander.Client(), firstTransferHash)
	waitForTxToBeIncludedInBatch(t, activeCommander.Client(), firstC2THash)
	waitForTxToBeIncludedInBatch(t, activeCommander.Client(), firstMMHash)
}
