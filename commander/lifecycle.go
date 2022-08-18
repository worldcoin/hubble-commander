package commander

import (
	"sync"
	"sync/atomic"
)

//nolint:structcheck
type lifecycle struct {
	mutex            sync.Mutex // protects Start method and startAndWaitChan
	startAndWaitChan chan struct{}

	active    uint32
	closeOnce sync.Once
}

func (l *lifecycle) isActive() bool {
	return atomic.LoadUint32(&l.active) != 0
}

func (l *lifecycle) setActive(active bool) {
	activeFlag := uint32(0)
	if active {
		activeFlag = 1
	}
	atomic.StoreUint32(&l.active, activeFlag)
}

func (l *lifecycle) getStartAndWaitChan() <-chan struct{} {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.unsafeGetStartAndWaitChan()
}

func (l *lifecycle) unsafeGetStartAndWaitChan() chan struct{} {
	if l.startAndWaitChan == nil {
		l.startAndWaitChan = make(chan struct{})
	}
	return l.startAndWaitChan
}

func (l *lifecycle) releaseStartAndWait() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	ch := l.unsafeGetStartAndWaitChan()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded by l.mutex.
		close(ch)
	}
}
