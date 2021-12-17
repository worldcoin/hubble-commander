package commander

import (
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) syncSpokes(startBlock, endBlock uint64) error {
	var newSpokesCount *int

	duration, err := metrics.MeasureDuration(func() (err error) {
		newSpokesCount, err = c.unmeasuredSyncSpokes(startBlock, endBlock)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.SyncingMethodDuration, prometheus.Labels{
		"method": metrics.SyncSpokesMethod,
	})

	logNewRegisteredSpokesCount(*newSpokesCount)

	return nil
}

func (c *Commander) unmeasuredSyncSpokes(startBlock, endBlock uint64) (*int, error) {
	newSpokesCount := 0

	it, err := c.getRegisteredSpokeIterator(startBlock, endBlock)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		if *isNewSpoke {
			newSpokesCount++
		}
	}

	return &newSpokesCount, it.Error()
}

func (c *Commander) getRegisteredSpokeIterator(start, end uint64) (*spokeregistry.RegisteredSpokeIterator, error) {
	it := &spokeregistry.RegisteredSpokeIterator{}

	err := c.client.FilterLogs(c.client.SpokeRegistry.BoundContract, "RegisteredSpoke", &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func saveSyncedSpoke(
	registeredSpokeStorage *st.RegisteredSpokeStorage,
	registeredSpoke *models.RegisteredSpoke,
) (isNewSpoke *bool, err error) {
	_, err = registeredSpokeStorage.GetRegisteredSpoke(registeredSpoke.ID)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	}

	if st.IsNotFoundError(err) {
		err = registeredSpokeStorage.AddRegisteredSpoke(registeredSpoke)
		if err != nil {
			return nil, err
		}
		return ref.Bool(true), nil
	} else {
		return ref.Bool(false), nil
	}
}

func logNewRegisteredSpokesCount(newSpokesCount int) {
	if newSpokesCount > 0 {
		log.Printf("Found %d new registered spoke(s)", newSpokesCount)
	}
}
