package commander

import (
	"context"
	"sync"
	"sync/atomic"
)

type lifecycle struct {
	active           uint32
	startAndWaitChan chan struct{}
	mutex            sync.Mutex
	manualStop       bool

	workersContext     context.Context
	stopWorkersContext context.CancelFunc
	workersWaitGroup   sync.WaitGroup
}

func (l *lifecycle) isActive() bool {
	return atomic.LoadUint32(&l.active) != 0
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

func (l *lifecycle) unsafeCloseStartAndWaitChan() {
	ch := l.unsafeGetStartAndWaitChan()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded by l.mutex.
		close(ch)
	}
}
