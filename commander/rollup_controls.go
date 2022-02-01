package commander

import (
	"context"
	"sync/atomic"
)

// nolint:structcheck
type rollupControls struct {
	batchCreationEnabled bool

	rollupLoopActive uint32
	cancelRollupLoop context.CancelFunc
}

func makeRollupControls() rollupControls {
	return rollupControls{
		batchCreationEnabled: true,
	}
}

func (c *rollupControls) isRollupLoopActive() bool {
	return atomic.LoadUint32(&c.rollupLoopActive) != 0
}

func (c *rollupControls) setRollupLoopActive(active bool) {
	activeFlag := uint32(0)
	if active {
		activeFlag = 1
	}
	atomic.StoreUint32(&c.rollupLoopActive, activeFlag)
}
