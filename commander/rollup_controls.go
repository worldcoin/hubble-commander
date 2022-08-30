package commander

import (
	"context"
	"sync/atomic"
)

//nolint:structcheck
type rollupControls struct {
	batchCreationEnabled bool
	migrate              uint32

	rollupLoopActive uint32
	cancelRollupLoop context.CancelFunc
}

func makeRollupControls(migrate bool) rollupControls {
	controls := rollupControls{
		batchCreationEnabled: true,
	}
	controls.setMigrate(migrate)
	return controls
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

func (c *rollupControls) isMigrating() bool {
	return atomic.LoadUint32(&c.migrate) != 0
}

func (c *rollupControls) setMigrate(migrate bool) {
	migrateFlag := uint32(0)
	if migrate {
		migrateFlag = 1
	}
	atomic.StoreUint32(&c.migrate, migrateFlag)
}
