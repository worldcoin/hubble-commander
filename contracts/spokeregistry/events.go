package spokeregistry

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type SpokeRegisteredIterator struct {
	SpokeRegistrySpokeRegisteredIterator
}

func (i *SpokeRegisteredIterator) SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription) {
	i.SpokeRegistrySpokeRegisteredIterator.contract = contract
	i.SpokeRegistrySpokeRegisteredIterator.event = event
	i.SpokeRegistrySpokeRegisteredIterator.logs = logs
	i.SpokeRegistrySpokeRegisteredIterator.sub = sub
}
