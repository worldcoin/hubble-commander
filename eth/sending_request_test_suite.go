package eth

import (
	"context"
	"sync"

	"github.com/stretchr/testify/require"
)

type testSuiteWithRequestsSending struct {
	senderCtxCancel func()
	wg              sync.WaitGroup
	requestsChan    chan *TxSendingRequest
}

func (ts *testSuiteWithRequestsSending) StartTxsSending(t require.TestingT, requests chan *TxSendingRequest) {
	var ctx context.Context
	ctx, ts.senderCtxCancel = context.WithCancel(context.Background())

	ts.wg.Add(1)
	go func() {
		err := ts.sendRequestsFromChan(ctx, requests)
		require.NoError(t, err)
		ts.wg.Done()
	}()
}

func (ts *testSuiteWithRequestsSending) StopTxsSending() {
	ts.senderCtxCancel()
	ts.wg.Wait()
}

func (ts *testSuiteWithRequestsSending) sendRequestsFromChan(ctx context.Context, requests <-chan *TxSendingRequest) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case req := <-requests:
			err := req.Send()
			if err != nil {
				return err
			}
		}
	}
}
