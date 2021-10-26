package chain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type RPCBackend struct {
	*ethclient.Client
}

func NewRPCBackend(c *rpc.Client) *RPCBackend {
	return &RPCBackend{
		Client: ethclient.NewClient(c),
	}
}

func (c *RPCBackend) Commit() {
	// NOOP
}
