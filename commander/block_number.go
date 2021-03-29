package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
)

var LatestBlockNumber uint32

func BlockNumberEndlessLoop(client *eth.Client, cfg *config.RollupConfig) error {
	done := make(chan bool)
	return BlockNumberLoop(client, cfg, done)
}

func BlockNumberLoop(client *eth.Client, cfg *config.RollupConfig, done <-chan bool) error {
	ticker := time.NewTicker(cfg.BlockNumberLoopInterval)

	for {
		select {
		case <-done:
			ticker.Stop()
			return nil
		case <-ticker.C:
			blockNumber, err := client.ChainConnection.GetLatestBlockNumber()
			if err != nil {
				log.Println(err.Error())
				return err
			}
			if *blockNumber > LatestBlockNumber {
				log.Printf("New block was mined: %d", blockNumber)
			}
			LatestBlockNumber = *blockNumber
		}
	}
}
