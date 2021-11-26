package depositmanager

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type DepositQueuedIterator struct {
	DepositManagerDepositQueuedIterator
}

func (i *DepositQueuedIterator) SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription) {
	i.DepositManagerDepositQueuedIterator.contract = contract
	i.DepositManagerDepositQueuedIterator.event = event
	i.DepositManagerDepositQueuedIterator.logs = logs
	i.DepositManagerDepositQueuedIterator.sub = sub
}

type DepositSubTreeReadyIterator struct {
	DepositManagerDepositSubTreeReadyIterator
}

func (i *DepositSubTreeReadyIterator) SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription) {
	i.DepositManagerDepositSubTreeReadyIterator.contract = contract
	i.DepositManagerDepositSubTreeReadyIterator.event = event
	i.DepositManagerDepositSubTreeReadyIterator.logs = logs
	i.DepositManagerDepositSubTreeReadyIterator.sub = sub
}
