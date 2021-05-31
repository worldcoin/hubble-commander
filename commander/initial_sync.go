package commander

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Commander) InitialSync() error {
	log.Println("Started initial syncing")
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}

	syncedBlock, err := c.storage.GetSyncedBlock(c.client.ChainState.ChainID)
	if err != nil {
		return err
	}
	startBlock := *syncedBlock
	endBlock := startBlock + uint64(c.cfg.Rollup.SyncSize)

	for endBlock != uint64(*latestBlockNumber) {
		//TODO: process accounts and batches parallel
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

		err = c.storage.SetSyncedBlock(c.client.ChainState.ChainID, endBlock)
		if err != nil {
			return err
		}

		startBlock = endBlock
		endBlock += uint64(c.cfg.Rollup.SyncSize)
		if endBlock > uint64(*latestBlockNumber) {
			endBlock = uint64(*latestBlockNumber)
		}
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
	for it.Next() {
		err = ProcessPubkeyRegistered(c.storage, it.Event)
		if err != nil {
			return err
		}
	}
	return nil
}
