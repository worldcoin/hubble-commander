package accountregistry

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type SinglePubKeyRegisteredIterator struct {
	AccountRegistrySinglePubkeyRegisteredIterator
}

func (i *SinglePubKeyRegisteredIterator) SetData(
	contract *bind.BoundContract,
	event string,
	logs chan types.Log,
	sub ethereum.Subscription,
) {
	i.AccountRegistrySinglePubkeyRegisteredIterator.contract = contract
	i.AccountRegistrySinglePubkeyRegisteredIterator.event = event
	i.AccountRegistrySinglePubkeyRegisteredIterator.logs = logs
	i.AccountRegistrySinglePubkeyRegisteredIterator.sub = sub
}

type BatchPubKeyRegisteredIterator struct {
	AccountRegistryBatchPubkeyRegisteredIterator
}

func (i *BatchPubKeyRegisteredIterator) SetData(
	contract *bind.BoundContract,
	event string,
	logs chan types.Log,
	sub ethereum.Subscription,
) {
	i.AccountRegistryBatchPubkeyRegisteredIterator.contract = contract
	i.AccountRegistryBatchPubkeyRegisteredIterator.event = event
	i.AccountRegistryBatchPubkeyRegisteredIterator.logs = logs
	i.AccountRegistryBatchPubkeyRegisteredIterator.sub = sub
}
