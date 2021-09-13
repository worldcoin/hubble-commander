package executor

import (
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type DisputeContext struct {
	storage *st.Storage
	client  *eth.Client
}

func NewDisputeContext(storage *st.Storage, client *eth.Client) *DisputeContext {
	return &DisputeContext{storage: storage, client: client}
}
