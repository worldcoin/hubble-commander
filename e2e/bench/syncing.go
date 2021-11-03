//go:build e2e
// +build e2e

package bench

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

func (s *BenchmarkSuite) benchSyncing() {
	cfg := config.GetConfig()

	cfg.Bootstrap.Prune = true
	cfg.API.Port = "5002"
	cfg.Badger.Path += "_passive"
	cfg.Bootstrap.ChainSpecPath = nil
	cfg.Bootstrap.BootstrapNodeURL = ref.String("http://localhost:8080")
	cfg.Ethereum.PrivateKey = "ab6919fd6ac00246bb78657e0696cf72058a4cb395133d074eabaddb83d8b00c"
	passiveCommander, err := setup.CreateInProcessCommanderWithConfig(cfg, false)
	s.NoError(err)
	err = passiveCommander.Start()
	s.NoError(err)
	defer func() {
		s.NoError(passiveCommander.Stop())
	}()

	// Observe commander syncing
	var networkInfo dto.NetworkInfo
	err = s.commander.Client().CallFor(&networkInfo, "hubble_getNetworkInfo")
	s.NoError(err)

	latestBatch := networkInfo.LatestBatch.Uint64()
	startTime := time.Now()
	lastSyncedBatch := uint64(0)
	for lastSyncedBatch < latestBatch {
		var networkInfo dto.NetworkInfo
		err = passiveCommander.Client().CallFor(&networkInfo, "hubble_getNetworkInfo")
		s.NoError(err)
		newBatch := uint64(0)
		if networkInfo.LatestBatch != nil {
			newBatch = networkInfo.LatestBatch.Uint64()
		}

		if newBatch == lastSyncedBatch {
			continue
		}
		lastSyncedBatch = newBatch

		txCount := networkInfo.TransactionCount

		fmt.Printf(
			"Transfers synced: %d, throughput: %f tx/s, batches synced: %d/%d\n",
			txCount,
			float64(txCount)/(time.Since(startTime).Seconds()),
			lastSyncedBatch,
			latestBatch,
		)
	}
}
