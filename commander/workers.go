package commander

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

type workers struct {
	workersContext     context.Context
	stopWorkersContext context.CancelFunc
	workersWaitGroup   sync.WaitGroup
}

func makeWorkers() workers {
	ctx, ctxCancel := context.WithCancel(context.Background())
	return workers{
		workersContext:     ctx,
		stopWorkersContext: ctxCancel,
	}
}

func (w *workers) startWorker(name string, fn func() error) {
	w.workersWaitGroup.Add(1)
	go func() {
		var err error
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				var ok bool
				err, ok = recoverErr.(error)
				if !ok {
					err = fmt.Errorf("%+v", recoverErr)
				}
			}
			if err != nil {
				log.Errorf("%s worker failed with: %+v", name, err)
				w.stopWorkersContext()
			}
			w.workersWaitGroup.Done()
		}()

		err = fn()
	}()
}

func (w *workers) stopWorkersAndWait() {
	w.stopWorkersContext()
	w.workersWaitGroup.Wait()
}
