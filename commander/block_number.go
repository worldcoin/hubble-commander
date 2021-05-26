package commander

import (
	"log"
	"time"
)

func (c *Commander) blockNumberLoop() error {
	ticker := time.NewTicker(c.cfg.Rollup.BlockNumberLoopInterval)

	for {
		select {
		case <-c.stopChannel:
			ticker.Stop()
			return nil
		case <-ticker.C:
			blockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
			if err != nil {
				log.Println(err.Error())
				return err
			}
			c.storage.SetLatestBlockNumber(*blockNumber)
		}
	}
}
