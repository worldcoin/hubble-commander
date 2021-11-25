package commander

import (
	"bytes"
	"context"

	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) syncTokens(startBlock, endBlock uint64) error {
	var newTokensCount *int

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
		tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.TokenRegistry.ABI.Methods["finaliseRegistration"].ID) {
			continue // TODO handle internal transactions
		}

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

	return &newTokensCount, nil
}

func (c *Commander) getRegisteredTokenIterator(start, end uint64) (*tokenregistry.RegisteredTokenIterator, error) {
	it := &tokenregistry.RegisteredTokenIterator{}

	err := c.client.FilterLogs(c.client.TokenRegistry.BoundContract, "RegisteredToken", &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func saveSyncedToken(
	registeredTokenStorage *st.RegisteredTokenStorage,
	registeredToken *models.RegisteredToken,
) (isNewToken *bool, err error) {
	_, err = registeredTokenStorage.GetRegisteredToken(registeredToken.ID)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	}

	if st.IsNotFoundError(err) {
		err = registeredTokenStorage.AddRegisteredToken(registeredToken)
		if err != nil {
			return nil, err
		}
		return ref.Bool(true), nil
	} else {
		return ref.Bool(false), nil
	}
}

func logNewRegisteredTokensCount(newTokensCount int) {
	if newTokensCount > 0 {
		log.Printf("Found %d new registered token(s)", newTokensCount)
	}
}
