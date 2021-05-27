package commander

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Commander) blockNumberLoop() error {
	//TODO: remove interval from cfg
	ticker := time.NewTicker(c.cfg.Rollup.BlockNumberLoopInterval)
	blocks := make(chan *types.Header)
	subscription, err := c.client.ChainConnection.SubscribeNewHead(blocks)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-c.stopChannel:
			ticker.Stop()
			return nil
		case <-blocks:
			//remove
			blockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
			if err != nil {
				log.Println(err.Error())
				return err
			}
			c.storage.SetLatestBlockNumber(*blockNumber)

			isProposer, err := c.client.IsActiveProposer()
			if err != nil {
				return errors.WithStack(err)
			}
			c.storage.SetProposer(isProposer)

			if !isProposer {
				err = c.SyncBatches()
				if err != nil {
					return err
				}
			}
		}
	}
}

func (c *Commander) SyncBatches() (err error) {
	transactionExecutor, err := newTransactionExecutor(c.storage, c.client, &c.cfg.Rollup)
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	return transactionExecutor.Commit()
}
