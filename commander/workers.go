package commander

import (
	"context"
	"sync"
)

type workers struct {
	workersContext     context.Context
	stopWorkersContext context.CancelFunc
	workersWaitGroup   sync.WaitGroup
}
