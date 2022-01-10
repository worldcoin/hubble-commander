package tokenregistry

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type RegisteredTokenIterator struct {
	TokenRegistryTokenRegisteredIterator
}

func (i *RegisteredTokenIterator) SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription) {
	i.TokenRegistryTokenRegisteredIterator.contract = contract
	i.TokenRegistryTokenRegisteredIterator.event = event
	i.TokenRegistryTokenRegisteredIterator.logs = logs
	i.TokenRegistryTokenRegisteredIterator.sub = sub
}
