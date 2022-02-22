package tracker

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/eth"
)

var errChannelClosed = fmt.Errorf("channel closed")

func (t *Tracker) SendRequestedTxs(ctx context.Context) error {
	err := t.sendRequestedTxs(ctx)
	close(t.requestsChan)

	for request := range t.requestsChan {
		request.ResultTxChan <- eth.SendResponse{Error: errChannelClosed}
	}
	return err
}

func (t *Tracker) sendRequestedTxs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case request := <-t.requestsChan:
			err := request.Send()
			if err != nil {
				return err
			}
		}
	}
}
