package commander

import (
	"bytes"
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) syncTokens(startBlock, endBlock uint64) error {
	startTime := time.Now()

	it, err := c.client.TokenRegistry.FilterRegisteredToken(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	})
	if err != nil {
		return err
	}
	defer func() { _ = it.Close() }()
	newTokensCount := 0

	for it.Next() {
		tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.TokenRegistryABI.Methods["finaliseRegistration"].ID) {
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
			return err
		}
		if *isNewToken {
			newTokensCount++
		}
	}

	measureRegisteredTokensSyncingDuration(startTime, c.metrics)
	logNewRegisteredTokensCount(newTokensCount)

	return nil
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

func measureRegisteredTokensSyncingDuration(
	start time.Time,
	commanderMetrics *metrics.CommanderMetrics,
) {
	duration := time.Since(start).Round(time.Millisecond)
	commanderMetrics.SyncingMethodDuration.WithLabelValues("sync_tokens").Observe(float64(duration.Milliseconds()))
}

func logNewRegisteredTokensCount(newTokensCount int) {
	if newTokensCount > 0 {
		log.Printf("Found %d new registered token(s)", newTokensCount)
	}
}
