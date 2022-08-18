package admin

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type API struct {
	cfg                 *config.APIConfig
	storage             *st.Storage
	client              *eth.Client
	enableBatchCreation func(enable bool)
	enableTxsAcceptance func(enable bool)
}

func NewAPI(
	cfg *config.APIConfig,
	storage *st.Storage,
	client *eth.Client,
	enableBatchCreation func(enable bool),
	enableTxsAcceptance func(enable bool),
) *API {
	return &API{
		cfg:                 cfg,
		storage:             storage,
		client:              client,
		enableBatchCreation: enableBatchCreation,
		enableTxsAcceptance: enableTxsAcceptance,
	}
}

func NewTestAPI(
	cfg *config.APIConfig,
	storage *st.Storage,
	client *eth.Client,
) *API {
	return NewAPI(
		cfg, storage, client, nil, nil,
	)
}
