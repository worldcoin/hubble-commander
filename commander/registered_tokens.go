package commander

import (
	"context"
	"errors"

	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	bh "github.com/timshannon/badgerhold/v4"
	"go.opentelemetry.io/otel"
)

func (c *Commander) syncTokens(ctx context.Context, startBlock, endBlock uint64) error {
	var newTokensCount *int

	_, span := otel.Tracer("rollupLoop").Start(ctx, "syncTokens")
	defer span.End()

	duration, err := metrics.MeasureDuration(func() (err error) {
		newTokensCount, err = c.unmeasuredSyncTokens(startBlock, endBlock)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.SyncingMethodDuration, prometheus.Labels{
		"method": metrics.SyncTokensMethod,
	})

	logNewRegisteredTokensCount(*newTokensCount)

	return nil
}

func (c *Commander) unmeasuredSyncTokens(startBlock, endBlock uint64) (*int, error) {
	newTokensCount := 0

	it, err := c.getRegisteredTokenIterator(startBlock, endBlock)
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	for it.Next() {
		tokenID := models.MakeUint256FromBig(*it.Event.TokenID)
		contract := it.Event.TokenContract
		registeredToken := &models.RegisteredToken{
			ID:       tokenID,
			Contract: contract,
		}

		isNewToken, err := saveSyncedToken(c.storage.RegisteredTokenStorage, registeredToken)
		if err != nil {
			return nil, err
		}
		if *isNewToken {
			newTokensCount++
		}
	}

	return &newTokensCount, it.Error()
}

func (c *Commander) getRegisteredTokenIterator(start, end uint64) (*tokenregistry.RegisteredTokenIterator, error) {
	it := &tokenregistry.RegisteredTokenIterator{}

	err := c.client.FilterLogs(c.client.TokenRegistry.BoundContract, eth.TokenRegisteredEvent, &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func saveSyncedToken(storage *st.RegisteredTokenStorage, token *models.RegisteredToken) (isNewToken *bool, err error) {
	err = storage.AddRegisteredToken(token)
	if errors.Is(err, bh.ErrKeyExists) {
		return ref.Bool(false), nil
	}
	if err != nil {
		return nil, err
	}
	return ref.Bool(true), nil
}

func logNewRegisteredTokensCount(newTokensCount int) {
	if newTokensCount > 0 {
		log.Printf("Found %d new registered token(s)", newTokensCount)
	}
}
