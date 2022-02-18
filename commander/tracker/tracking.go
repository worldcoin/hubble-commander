package tracker

import (
	"context"
	"fmt"
	"sync"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func TrackSentTxs(ctx context.Context, client *eth.Client, sentTxsChan <-chan *types.Transaction) error {
	queue := newTxsQueue()
	return trackSentTxs(ctx, client, sentTxsChan, queue)
}

func trackSentTxs(
	ctx context.Context,
	client *eth.Client,
	sentTxsChan <-chan *types.Transaction,
	queue *txsQueue,
) error {
	wg := sync.WaitGroup{}
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg.Add(1)
	go func(queue *txsQueue) {
		startReadingChannel(subCtx, queue, sentTxsChan)
		wg.Done()
	}(queue)

	errChan := make(chan error)
	wg.Add(1)
	go func(queue *txsQueue) {
		err := startCheckingTxs(subCtx, queue, client)
		errChan <- err
		wg.Done()
	}(queue)

	err := <-errChan
	cancel()
	wg.Wait()
	return err
}

func startReadingChannel(ctx context.Context, queue *txsQueue, txsChan <-chan *types.Transaction) {
	for {
		select {
		case <-ctx.Done():
			return
		case tx := <-txsChan:
			queue.Add(tx)
		}
	}
}

func startCheckingTxs(ctx context.Context, queue *txsQueue, client *eth.Client) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if queue.IsEmpty() {
				continue
			}
			tx := queue.First()
			err := waitUntilTxMinedAndCheckForFail(client, tx)
			if err != nil {
				return err
			}

			queue.RemoveFirst()
		}
	}
}

func waitUntilTxMinedAndCheckForFail(client *eth.Client, tx *types.Transaction) error {
	receipt, err := client.WaitToBeMined(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if receipt.Status == 1 {
		return nil
	}
	err = client.GetRevertMessage(tx, receipt)
	return fmt.Errorf("%w tx_hash=%s", err, tx.Hash().String())
}
