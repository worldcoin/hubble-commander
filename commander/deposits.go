package commander

import (
	"context"

	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *Commander) syncDeposits(ctx context.Context, start, end uint64) error {
	var depositSubtrees []models.PendingDepositSubtree

	_, span := rollupTracer.Start(ctx, "syncDeposits")
	defer span.End()

	duration, err := metrics.MeasureDuration(func() error {
		var err error

		err = c.syncQueuedDeposits(start, end)
		if err != nil {
			return err
		}

		depositSubtrees, err = c.fetchDepositSubtrees(start, end)
		if err != nil {
			return err
		}

		if len(depositSubtrees) > 0 {
			return c.saveSyncedSubtrees(depositSubtrees)
		}

		return nil
	})
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.SyncingMethodDuration, prometheus.Labels{
		"method": metrics.SyncDepositsMethod,
	})

	return nil
}

func (c *Commander) syncQueuedDeposits(start, end uint64) error {
	it, err := c.getDepositQueuedIterator(start, end)
	if err != nil {
		return err
	}
	defer func() { _ = it.Close() }()

	for it.Next() {
		deposit := models.PendingDeposit{
			ID: models.DepositID{
				SubtreeID:    models.MakeUint256FromBig(*it.Event.SubtreeID),
				DepositIndex: models.MakeUint256FromBig(*it.Event.DepositID),
			},
			ToPubKeyID: uint32(it.Event.PubkeyID.Uint64()),
			TokenID:    models.MakeUint256FromBig(*it.Event.TokenID),
			L2Amount:   models.MakeUint256FromBig(*it.Event.L2Amount),
		}

		err = c.storage.AddPendingDeposit(&deposit)
		if err != nil {
			return err
		}
	}

	return it.Error()
}

func (c *Commander) fetchDepositSubtrees(start, end uint64) ([]models.PendingDepositSubtree, error) {
	it, err := c.getDepositSubtreeReadyIterator(start, end)
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	depositSubtrees := make([]models.PendingDepositSubtree, 0, 1)

	for it.Next() {
		subtree := models.PendingDepositSubtree{
			ID:   models.MakeUint256FromBig(*it.Event.SubtreeID),
			Root: it.Event.SubtreeRoot,
		}

		depositSubtrees = append(depositSubtrees, subtree)
	}

	return depositSubtrees, it.Error()
}

func (c *Commander) getDepositQueuedIterator(start, end uint64) (*depositmanager.DepositQueuedIterator, error) {
	it := &depositmanager.DepositQueuedIterator{}

	err := c.client.FilterLogs(c.client.DepositManager.BoundContract, eth.DepositQueuedEvent, &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func (c *Commander) getDepositSubtreeReadyIterator(start, end uint64) (*depositmanager.DepositSubTreeReadyIterator, error) {
	it := &depositmanager.DepositSubTreeReadyIterator{}

	err := c.client.FilterLogs(c.client.DepositManager.BoundContract, eth.DepositSubTreeReadyEvent, &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func (c *Commander) saveSyncedSubtrees(subtrees []models.PendingDepositSubtree) error {
	maxDepositSubtreeDepth, err := c.client.GetMaxSubtreeDepthParam()
	if err != nil {
		return err
	}

	subtreeLeavesAmount := 1 << *maxDepositSubtreeDepth

	for i := range subtrees {
		err := c.saveSingleSubtree(&subtrees[i], subtreeLeavesAmount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commander) saveSingleSubtree(subtree *models.PendingDepositSubtree, subtreeLeavesAmount int) error {
	return c.storage.ExecuteInTransaction(st.TxOptions{}, func(txStorage *st.Storage) error {
		deposits, err := txStorage.GetFirstPendingDeposits(subtreeLeavesAmount)
		if err != nil {
			return err
		}

		subtree.Deposits = deposits

		err = txStorage.AddPendingDepositSubtree(subtree)
		if err != nil {
			return err
		}

		return txStorage.RemovePendingDeposits(deposits)
	})
}
