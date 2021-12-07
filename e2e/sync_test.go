//go:build e2e
// +build e2e

package e2e

import (
	"os"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/testutils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCommanderSync(t *testing.T) {
	cfg := config.GetConfig()
	if cfg.Ethereum == nil {
		log.Panicf("sync test cannot be run on simulator")
	}

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5001"
	cfg.Metrics.Port = "2001"
	activeCommander, err := setup.CreateInProcessCommanderWithConfig(cfg, true)
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
	testSubmitTransferBatch(t, activeCommander.Client(), senderWallet, 0)

	firstC2TWallet := wallets[len(wallets)-32]
	testSubmitC2TBatch(t, activeCommander.Client(), senderWallet, wallets, firstC2TWallet.PublicKey(), 32)

	testSubmitDepositBatch(t, activeCommander.Client())

	testSubmitMassMigrationBatch(t, activeCommander.Client(), senderWallet, 64)

	waitForBatch(t, activeCommander.Client(), models.MakeUint256(4))

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5002"
	cfg.Metrics.Port = "2002"
	cfg.Badger.Path += "_passive"
	cfg.Ethereum.PrivateKey = "ab6919fd6ac00246bb78657e0696cf72058a4cb395133d074eabaddb83d8b00c"
	passiveCommander, err := setup.CreateInProcessCommanderWithConfig(cfg, false)
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

	testSenderStateAfterTransfers(t, passiveCommander.Client(), senderWallet)
	testFeeReceiverStateAfterTransfers(t, passiveCommander.Client(), feeReceiverWallet)
	testGetBatches(t, passiveCommander.Client())
}
