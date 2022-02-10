package tracker

import (
	"context"
	"sync"

	"github.com/Worldcoin/hubble-commander/eth"
)

type TestSuiteWithTxsSending struct {
	txsChannel    chan *eth.TxSendingRequest
	cancelSending func()
	waitGroup     sync.WaitGroup
}

func (t *TestSuiteWithTxsSending) initTxsSending(channel chan *eth.TxSendingRequest) {
	if channel == nil {
		channel = make(chan *eth.TxSendingRequest, 32)
	}

	t.txsChannel = channel
}

func (t *TestSuiteWithTxsSending) StartTxsSending(channel chan *eth.TxSendingRequest) {
	t.initTxsSending(channel)
	t.waitGroup.Add(1)

	var ctx context.Context
	ctx, t.cancelSending = context.WithCancel(context.Background())
	go func() {
		_ = StartTxsRequestsSending(ctx, t.txsChannel)
		t.waitGroup.Done()
	}()
}

func (t *TestSuiteWithTxsSending) StopTxsSending() {
	t.cancelSending()
	t.waitGroup.Wait()
}
