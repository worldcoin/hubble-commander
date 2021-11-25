package tokenregistry

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type RegisteredTokenIterator struct {
	TokenRegistryRegisteredTokenIterator
}

func (i *RegisteredTokenIterator) SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription) {
	i.TokenRegistryRegisteredTokenIterator.contract = contract
	i.TokenRegistryRegisteredTokenIterator.event = event
	i.TokenRegistryRegisteredTokenIterator.logs = logs
	i.TokenRegistryRegisteredTokenIterator.sub = sub
}
