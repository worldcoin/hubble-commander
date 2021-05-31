package commander

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Commander) newBlockLoop() error {
	err := c.SyncOnStart()
	if err != nil {
		return err
	}

	blocks := make(chan *types.Header)
	subscription, err := c.client.ChainConnection.SubscribeNewHead(blocks)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	isProposer, err := c.client.IsActiveProposer()
	if err != nil {
		return errors.WithStack(err)
	}
	c.storage.SetProposer(isProposer)

	for {
		select {
		case <-c.stopChannel:
			return nil
		case err = <-subscription.Err():
			return err
		case newBlock := <-blocks:
			c.storage.SetLatestBlockNumber(uint32(newBlock.Number.Uint64()))

			isProposer, err = c.client.IsActiveProposer()
			if err != nil {
				return errors.WithStack(err)
			}
			c.storage.SetProposer(isProposer)

			err = c.SyncBatches(isProposer, nil)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
}

func (c *Commander) SyncBatches(isProposer bool, endBlock *uint64) (err error) {
	if isProposer {
		return nil
	}

	transactionExecutor, err := newTransactionExecutor(c.storage, c.client, c.cfg.Rollup)
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	err = transactionExecutor.SyncBatches(endBlock)
	if err != nil {
		return err
	}
	return transactionExecutor.Commit()
}

func (c *Commander) SyncOnStart() error {
	log.Println("Started initial syncing")
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}

	syncedBlock, err := c.storage.GetSyncedBlock(c.client.ChainState.ChainID)
	if err != nil {
		return err
	}
	startBlock := uint64(*syncedBlock)
	endBlock := startBlock + uint64(c.cfg.Rollup.SyncSize)

	for endBlock <= uint64(*latestBlockNumber) {
		err = c.RegisterAccounts(&bind.FilterOpts{
			Start: startBlock,
			End:   &endBlock,
		})
		if err != nil {
			return err
		}

		err = c.SyncBatches(false, &endBlock)
		if err != nil {
			return err
		}

		err = c.storage.SetSyncedBlock(c.client.ChainState.ChainID, uint32(endBlock))
		if err != nil {
			return err
		}

		startBlock = endBlock
		endBlock += uint64(c.cfg.Rollup.SyncSize)
	}

	log.Println("Finished initial syncing")
	return nil
}

func (c *Commander) RegisterAccounts(opts *bind.FilterOpts) error {
	it, err := c.client.AccountRegistry.FilterPubkeyRegistered(opts)
	if err != nil {
		return err
	}
	defer it.Close()
	i := 0
	for it.Next() {
		fmt.Printf("i: %d\n", i)
		i++
		err = ProcessPubkeyRegistered(c.storage, it.Event)
		if err != nil {
			return err
		}
	}
	return nil
}
