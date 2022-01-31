package admin

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type API struct {
	cfg     *config.APIConfig
	storage *st.Storage
	client  *eth.Client
}

func NewAPI(cfg *config.APIConfig, storage *st.Storage, client *eth.Client) *API {
	return &API{cfg: cfg, storage: storage, client: client}
}
