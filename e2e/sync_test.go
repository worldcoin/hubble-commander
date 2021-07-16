// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCommanderSync(t *testing.T) {
	cfg := config.GetConfig()
	if cfg.Ethereum == nil {
		log.Fatal("sync test cannot be run on simulator")
	}

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5001"
	activeCommander := setup.CreateInProcessCommanderWithConfig(cfg)

	err := activeCommander.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, activeCommander.Stop())
	}()

	domain := getDomain(t, activeCommander.Client())

	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	feeReceiverWallet := wallets[0]
	senderWallet := wallets[1]

	testGetVersion(t, activeCommander.Client())
	testGetUserStates(t, activeCommander.Client(), senderWallet)
	firstTransferHash := testSendTransfer(t, activeCommander.Client(), senderWallet, models.NewUint256(0))
	testGetTransaction(t, activeCommander.Client(), firstTransferHash)
	send31MoreTransfers(t, activeCommander.Client(), senderWallet, 1)

	firstC2TWallet := wallets[len(wallets)-32]
	firstCreate2TransferHash := testSendCreate2Transfer(t, activeCommander.Client(), senderWallet, *firstC2TWallet.PublicKey())
	testGetTransaction(t, activeCommander.Client(), firstCreate2TransferHash)
	send31MoreCreate2Transfers(t, activeCommander.Client(), senderWallet, wallets)

	waitForTxToBeIncludedInBatch(t, activeCommander.Client(), firstTransferHash)
	waitForTxToBeIncludedInBatch(t, activeCommander.Client(), firstCreate2TransferHash)

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5002"
	cfg.Badger.Path += "_passive"
	cfg.Postgres.Name += "_passive"
	cfg.Bootstrap.BootstrapNodeURL = ref.String("http://localhost:5001")
	cfg.Ethereum.PrivateKey = "ab6919fd6ac00246bb78657e0696cf72058a4cb395133d074eabaddb83d8b00c"
	passiveCommander := setup.CreateInProcessCommanderWithConfig(cfg)

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
