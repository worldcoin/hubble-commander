package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) syncTokens(startBlock, endBlock uint64) error {
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

	logRegisteredTokensCount(newTokensCount)
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

func logRegisteredTokensCount(newTokensCount int) {
	if newTokensCount > 0 {
		log.Printf("Found %d new registered token(s)", newTokensCount)
	}
}
