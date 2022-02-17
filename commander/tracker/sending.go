package tracker

import (
	"context"

	"github.com/Worldcoin/hubble-commander/eth"
)

func StartTxsRequestsSending(ctx context.Context, requestsChan <-chan *eth.TxSendingRequest) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case request := <-requestsChan:
			err := request.Send()
			if err != nil {
				return err
			}
		}
	}
}
