package commander

import (
	"context"
	"sync"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

// TODO-dis: replace with uint64 as it's not shared between workers
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
	if i.Get() == 0 {
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

//TODO: remove
func (c *Commander) watchDisputes() error {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := c.client.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	if err != nil {
		return errors.WithStack(err)
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-c.workersContext.Done():
			return nil
		case err = <-subscription.Err():
			return errors.WithStack(err)
		case rollbackStatus := <-sink:
			if rollbackStatus.Completed {
				continue
			}

			_, err = c.client.KeepRollingBack()
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
}

func (c *Commander) handleBatchRollback(rollupCancel context.CancelFunc) (bool, error) {
	invalidBatchID, err := c.client.GetInvalidBatchID()
	if err != nil {
		return false, err
	}
	if invalidBatchID == 0 {
		c.invalidBatchID.Reset()
		return false, nil
	}

	return true, c.manageRemoteBatchRollback(invalidBatchID, rollupCancel)
}

func (c *Commander) manageRemoteBatchRollback(batchID uint64, rollupCancel context.CancelFunc) error {
	c.invalidBatchID.Set(batchID)
	c.stopRollupLoop(rollupCancel)
	// TODO-dis trigger keep rolling back and handle errors (2)

	invalidBatch, err := c.storage.GetBatch(models.MakeUint256(batchID))
	if st.IsNotFoundError(err) {
		return nil
	}
	if err != nil {
		return err
	}

	return c.revertBatches(invalidBatch)
}
