package spokeregistry

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type RegisteredSpokeIterator struct {
	SpokeRegistryRegisteredSpokeIterator
}

func (i *RegisteredSpokeIterator) SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription) {
	i.SpokeRegistryRegisteredSpokeIterator.contract = contract
	i.SpokeRegistryRegisteredSpokeIterator.event = event
	i.SpokeRegistryRegisteredSpokeIterator.logs = logs
	i.SpokeRegistryRegisteredSpokeIterator.sub = sub
}
