package applier

import (
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type Applier struct {
	storage *st.Storage
	client  *eth.Client
}

func NewApplier(storage *st.Storage, client *eth.Client) *Applier {
	return &Applier{
		storage: storage,
		client:  client,
	}
}
