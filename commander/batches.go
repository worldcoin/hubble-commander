package commander

import (
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func BatchesEndlessLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	done := make(chan bool)
	return BatchesLoop(storage, client, cfg, done)
}

func BatchesLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, done <-chan bool) error {
	ticker := time.NewTicker(cfg.BatchLoopInterval)

	for {
		select {
		case <-done:
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := SubmitTransactionBatch(storage, client, cfg)
			if err != nil {
				return err
			}
		}
	}
}

func SubmitTransactionBatch(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	// _, err = client.SubmitTransfersBatch([]*models.Commitment{commitment})
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Sumbmited commitment %s on chain", commitment.LeafHash().Hex())
	return nil
}
