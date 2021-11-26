package rollup

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type NewBatchIterator struct {
	RollupNewBatchIterator
}

func (i *NewBatchIterator) SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription) {
	i.RollupNewBatchIterator.contract = contract
	i.RollupNewBatchIterator.event = event
	i.RollupNewBatchIterator.logs = logs
	i.RollupNewBatchIterator.sub = sub
}

type DepositsFinalisedIterator struct {
	RollupDepositsFinalisedIterator
}

func (i *DepositsFinalisedIterator) SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription) {
	i.RollupDepositsFinalisedIterator.contract = contract
	i.RollupDepositsFinalisedIterator.event = event
	i.RollupDepositsFinalisedIterator.logs = logs
	i.RollupDepositsFinalisedIterator.sub = sub
}
