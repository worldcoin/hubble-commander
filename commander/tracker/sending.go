package tracker

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/ethereum/go-ethereum/core/types"
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
			tx, err := t.sendTx(request)
			if err != nil {
				return err
			}
			t.txsChan <- tx
		}
	}
}

func (t *Tracker) sendTx(request *eth.TxSendingRequest) (*types.Transaction, error) {
	tx, err := request.Send(t.nonce)
	if err != nil {
		return nil, err
	}
	t.nonce++
	return tx, nil
}
