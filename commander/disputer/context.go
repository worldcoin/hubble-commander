package disputer

import (
	"github.com/Worldcoin/hubble-commander/commander/prover"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type Context struct {
	storage   *st.Storage
	client    *eth.Client
	proverCtx *prover.Context
}

func NewContext(storage *st.Storage, client *eth.Client) *Context {
	return &Context{
		storage:   storage,
		client:    client,
		proverCtx: prover.NewContext(storage),
	}
}
