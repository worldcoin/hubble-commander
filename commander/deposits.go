package commander

import (
	"bytes"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Commander) syncDeposits(start, end uint64) error {
	var depositSubTrees []models.PendingDepositSubTree

	syncDepositsAndFetchSubTreesDuration, err := metrics.MeasureDuration(func() error {
		var err error

		err = c.syncQueuedDeposits(start, end)
		if err != nil {
			return err
		}

		depositSubTrees, err = c.fetchDepositSubTrees(start, end)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	if len(depositSubTrees) > 0 {
		syncSubTreesDuration, err := metrics.MeasureDuration(func() error {
			return c.saveSyncedSubTrees(depositSubTrees)
		})
		if err != nil {
			return err
		}

		totalDuration := time.Duration(syncDepositsAndFetchSubTreesDuration.Milliseconds() + syncSubTreesDuration.Milliseconds())
		saveSyncDepositsWithNewSubTreesDurationMeasurement(totalDuration, c.metrics)
	} else {
		saveSyncDepositsWithNoNewSubTreesDurationMeasurement(*syncDepositsAndFetchSubTreesDuration, c.metrics)
	}

	return nil
}

func (c *Commander) syncQueuedDeposits(start, end uint64) error {
	it, err := c.client.DepositManager.FilterDepositQueued(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return err
	}
	defer func() { _ = it.Close() }()

	for it.Next() {
		tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.DepositManagerABI.Methods["depositFor"].ID) {
			continue // TODO handle internal transactions
		}

		deposit := models.PendingDeposit{
			ID: models.DepositID{
				BlockNumber: uint32(it.Event.Raw.BlockNumber),
				LogIndex:    uint32(it.Event.Raw.Index),
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

	return nil
}

func (c *Commander) fetchDepositSubTrees(start, end uint64) ([]models.PendingDepositSubTree, error) {
	it, err := c.client.DepositManager.FilterDepositSubTreeReady(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	depositSubTrees := make([]models.PendingDepositSubTree, 0, 1)

	for it.Next() {
		tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.DepositManagerABI.Methods["depositFor"].ID) {
			continue // TODO handle internal transactions
		}

		subTree := models.PendingDepositSubTree{
			ID:   models.MakeUint256FromBig(*it.Event.SubtreeID),
			Root: it.Event.SubtreeRoot,
		}

		depositSubTrees = append(depositSubTrees, subTree)
	}

	return depositSubTrees, nil
}

func (c *Commander) saveSyncedSubTrees(subTrees []models.PendingDepositSubTree) error {
	maxDepositSubTreeDepth, err := c.client.GetMaxSubTreeDepthParam()
	if err != nil {
		return err
	}

	subTreeLeavesAmount := 1 << *maxDepositSubTreeDepth

	for i := range subTrees {
		err := c.saveSingleSubTree(&subTrees[i], subTreeLeavesAmount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commander) saveSingleSubTree(subTree *models.PendingDepositSubTree, subTreeLeavesAmount int) error {
	return c.storage.ExecuteInTransaction(st.TxOptions{}, func(txStorage *st.Storage) error {
		deposits, err := txStorage.GetFirstPendingDeposits(subTreeLeavesAmount)
		if err != nil {
			return err
		}

		subTree.Deposits = deposits

		err = txStorage.AddPendingDepositSubTree(subTree)
		if err != nil {
			return err
		}

		return txStorage.RemovePendingDeposits(deposits)
	})
}

func saveSyncDepositsWithNoNewSubTreesDurationMeasurement(
	duration time.Duration,
	commanderMetrics *metrics.CommanderMetrics,
) {
	commanderMetrics.SyncingMethodDuration.
		With(prometheus.Labels{
			"method": metrics.SyncDepositsWithNoNewSubTreesMethod,
		}).
		Observe(float64(duration.Milliseconds()))
}

func saveSyncDepositsWithNewSubTreesDurationMeasurement(
	duration time.Duration,
	commanderMetrics *metrics.CommanderMetrics,
) {
	commanderMetrics.SyncingMethodDuration.
		With(prometheus.Labels{
			"method": metrics.SyncDepositsWithNewSubTreesMethod,
		}).
		Observe(float64(duration.Milliseconds()))
}
