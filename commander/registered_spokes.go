package commander

import (
	"context"
	"errors"

	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	bh "github.com/timshannon/badgerhold/v4"
)

func (c *Commander) syncSpokes(ctx context.Context, startBlock, endBlock uint64) error {
	_, span := newBlockTracer.Start(ctx, "syncSpokes")
	defer span.End()

	duration, err := metrics.MeasureDuration(func() error {
		return c.unmeasuredSyncSpokes(startBlock, endBlock)
	})
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.SyncingMethodDuration, prometheus.Labels{
		"method": metrics.SyncSpokesMethod,
	})

	return nil
}

func (c *Commander) unmeasuredSyncSpokes(startBlock, endBlock uint64) error {
	newSpokesCount := 0

	it, err := c.getSpokeRegisteredIterator(startBlock, endBlock)
	if err != nil {
		return err
	}
	defer func() { _ = it.Close() }()

	for it.Next() {
		spokeID := models.MakeUint256FromBig(*it.Event.SpokeID)
		contract := it.Event.SpokeContract
		registeredSpoke := &models.RegisteredSpoke{
			ID:       spokeID,
			Contract: contract,
		}

		isNewSpoke, err := saveSyncedSpoke(c.storage.RegisteredSpokeStorage, registeredSpoke)
		if err != nil {
			return err
		}
		if *isNewSpoke {
			newSpokesCount++
		}
	}
	if it.Error() != nil {
		return it.Error()
	}

	logNewRegisteredSpokesCount(newSpokesCount)
	return nil
}

func (c *Commander) getSpokeRegisteredIterator(start, end uint64) (*spokeregistry.SpokeRegisteredIterator, error) {
	it := &spokeregistry.SpokeRegisteredIterator{}

	err := c.client.FilterLogs(c.client.SpokeRegistry.BoundContract, eth.SpokeRegisteredEvent, &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func saveSyncedSpoke(storage *st.RegisteredSpokeStorage, spoke *models.RegisteredSpoke) (isNewSpoke *bool, err error) {
	err = storage.AddRegisteredSpoke(spoke)
	if errors.Is(err, bh.ErrKeyExists) {
		return ref.Bool(false), nil
	}
	if err != nil {
		return nil, err
	}
	return ref.Bool(true), nil
}

func logNewRegisteredSpokesCount(newSpokesCount int) {
	if newSpokesCount > 0 {
		log.Printf("Found %d new registered spoke(s)", newSpokesCount)
	}
}
