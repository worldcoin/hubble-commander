package commander

import (
	"context"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *Commander) syncStakeWithdrawals(ctx context.Context, startBlock, endBlock uint64) error {

	_, span := rollupTracer.Start(ctx, "syncStakeWithdrawls")
	defer span.End()

	duration, err := metrics.MeasureDuration(func() (err error) {
		err = c.unmeasuredSyncStakeWithdrawals(startBlock, endBlock)
		return err
	})
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.SyncingMethodDuration, prometheus.Labels{
		"method": metrics.SyncStakeWithdrawalsMethod,
	})
	return nil
}

func (c *Commander) unmeasuredSyncStakeWithdrawals(startBlock, endBlock uint64) error {
	it, err := c.getStakeWithdrawIterator(startBlock, endBlock)
	if err != nil {
		return err
	}
	defer func() { _ = it.Close() }()

	for it.Next() {
		if it.Event.Committed != c.client.Blockchain.GetAccount().From {
			continue
		}

		err = c.storage.RemovePendingStakeWithdrawal(models.MakeUint256FromBig(*it.Event.BatchID))
		if err != nil && !storage.IsNotFoundError(err) {
			return err
		}
	}
	return nil
}

func (c *Commander) getStakeWithdrawIterator(start, end uint64) (*rollup.StakeWithdrawIterator, error) {
	it := &rollup.StakeWithdrawIterator{}

	err := c.client.FilterLogs(c.client.Rollup.BoundContract, eth.StakeWithdrawEvent, &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}
