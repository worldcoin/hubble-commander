package syncer

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type DepositsContext struct {
	storage *st.Storage
	client  *eth.Client
	applier *applier.Applier
}

func newDepositsContext(storage *st.Storage, client *eth.Client) *DepositsContext {
	return &DepositsContext{
		storage: storage,
		client:  client,
		applier: applier.NewApplier(storage),
	}
}

func (c *DepositsContext) Commit() {}

func (c *DepositsContext) Rollback() {}
