package tracker

import (
	"context"
)

func (t *Tracker) SendRequestedTxs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case request := <-t.requestsChan:
			err := request.Send()
			if err != nil {
				// nolint:gocritic
				// close(t.requestsChan)
				return err
			}
		}
	}
}
