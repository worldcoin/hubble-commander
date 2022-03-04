package tracker

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/eth"
)

var errChannelClosed = fmt.Errorf("channel closed")

func (t *Tracker) sendRequestedTxs(ctx context.Context) error {
	err := t.sendRequestedTxsLoop(ctx)
	close(t.requestsChan)

	for request := range t.requestsChan {
		request.ResultTxChan <- eth.SendResponse{Error: errChannelClosed}
	}
	return err
}

func (t *Tracker) sendRequestedTxsLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case request := <-t.requestsChan:
			if err := t.sendTx(request); err != nil {
				return err
			}
		}
	}
}

func (t *Tracker) sendTx(request *eth.TxSendingRequest) error {
	tx, err := request.Send(t.nonce)
	if err != nil {
		return err
	}
	t.nonce++
	if request.ShouldTrackTx {
		t.txsChan <- tx
	}
	return nil
}
