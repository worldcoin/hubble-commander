package syncer

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type DepositsContext struct {
	cfg     *config.RollupConfig
	storage *st.Storage
	client  *eth.Client
	applier *applier.Applier
}

func NewTestDepositsContext(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) *DepositsContext {
	return newDepositsContext(storage, client, cfg)
}

func newDepositsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
) *DepositsContext {
	return &DepositsContext{
		cfg:     cfg,
		storage: storage,
		client:  client,
		applier: applier.NewApplier(storage, client),
	}
}
