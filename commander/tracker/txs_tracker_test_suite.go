package tracker

import (
	"context"
	"sync"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

type TestSuiteWithTxsTracker struct {
	TxsTracker    *TxsTracker
	cancelTracker func()
	waitGroup     sync.WaitGroup
}

func (th *TestSuiteWithTxsTracker) InitTracker(client *eth.Client, txsChan chan *types.Transaction) {
	if txsChan == nil {
		txsChan = make(chan *types.Transaction, 32)
	}
	th.TxsTracker = NewTxTracker(client, txsChan)
	th.waitGroup = sync.WaitGroup{}
}

func (th *TestSuiteWithTxsTracker) StartTracker(t require.TestingT) {
	th.waitGroup.Add(1)

	var ctx context.Context
	ctx, th.cancelTracker = context.WithCancel(context.Background())
	go func() {
		defer th.waitGroup.Done()
		err := th.TxsTracker.StartTracking(ctx)
		require.NoError(t, err)
	}()
}

func (th *TestSuiteWithTxsTracker) StopTracker() {
	th.cancelTracker()
	th.waitGroup.Wait()
}
