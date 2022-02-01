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

func (c *workers) startWorker(name string, fn func() error) {
	c.workersWaitGroup.Add(1)
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
				c.stopWorkersContext()
			}
			c.workersWaitGroup.Done()
		}()

		err = fn()
	}()
}
