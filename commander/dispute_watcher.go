package commander

import (
	"sync"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

type InvalidBatchID struct {
	id    uint64
	mutex sync.RWMutex
}

func (i *InvalidBatchID) Set(id uint64) {
	if id == 0 {
		return
	}

	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.id = id
}

func (i *InvalidBatchID) Reset() {
	//TODO-dis: should be read lock here
	if i.id == 0 {
		return
	}
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.id = 0
}

func (i *InvalidBatchID) Get() uint64 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.id
}

/*
	invalidBatchID := getInvalidBatchMarker
	if invalidBatchID != 0 {
		stop rollup loop if running

		lastBatchID := local last batch id
		if lastBatchID >= invalidBatchID {
			revertBatchesInRange
		}
		sync accounts normally
		sync batches only to invalidBatchID - 1

	}
*/

func (c *Commander) watchDisputes() error {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := c.client.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err = <-subscription.Err():
			return errors.WithStack(err)
		case rollbackStatus := <-sink:
			if rollbackStatus.Completed {
				c.invalidBatchID.Reset()
			}
			//TODO-dis: call keep rolling back in case it's not completed
			//transactionHash, err = c.client.Rollup.keepRollingBack()
			//if err != nil {
			//	return err
			//}
		}
	}
}
