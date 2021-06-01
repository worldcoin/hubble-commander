package commander

import (
	"context"
	"log"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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

	group, ctx := errgroup.WithContext(context.Background())
	boundaries := make(chan uint64, 5)

	group.Go(func() error {
		return c.initialSyncAccounts(ctx, boundaries, *syncedBlock, uint64(*latestBlockNumber))
	})
	group.Go(func() error {
		return c.initialSyncBatches(ctx, boundaries)
	})
	err = group.Wait()
	if err != nil {
		return errors.WithStack(err)
	}

	log.Println("Finished initial syncing")
	return nil
}

func (c *Commander) initialSyncAccounts(ctx context.Context, boundaries chan<- uint64, startBlock, latestBlock uint64) error {
	defer close(boundaries)
	opts := &bind.FilterOpts{
		Start: startBlock,
		End:   ref.Uint64(startBlock + uint64(c.cfg.Rollup.SyncSize)),
	}
	for *opts.End != latestBlock {
		err := c.registerAccounts(opts)
		if err != nil {
			return err
		}
		opts.Start = *opts.End
		*opts.End = calculateEndBlock(*opts.End, latestBlock, c.cfg.Rollup.SyncSize)
		boundaries <- *opts.End
		select {
		case <-ctx.Done():
			return nil
		default:
			continue
		}
	}
	return nil
}

func (c *Commander) initialSyncBatches(ctx context.Context, boundaries <-chan uint64) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case endBlock, ok := <-boundaries:
			if !ok {
				return nil
			}
			err := c.SyncBatches(false, &endBlock)
			if err != nil {
				return err
			}
		}
	}
}

func calculateEndBlock(endBlock, latestBlock uint64, syncSize uint32) uint64 {
	endBlock += uint64(syncSize)
	if endBlock > latestBlock {
		return latestBlock
	}
	return endBlock
}

func (c *Commander) registerAccounts(opts *bind.FilterOpts) error {
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
