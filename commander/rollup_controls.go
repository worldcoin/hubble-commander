package commander

import "context"

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
