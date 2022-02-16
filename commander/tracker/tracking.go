package tracker

import (
	"context"
	"fmt"
	"sync"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func StartTrackingSentTxs(ctx context.Context, client *eth.Client, sentTxsChan <-chan *types.Transaction) error {
	queue := newTxsQueue()
	return startTrackingSentTxs(ctx, client, sentTxsChan, queue)
}

func startTrackingSentTxs(
	ctx context.Context,
	client *eth.Client,
	sentTxsChan <-chan *types.Transaction,
	queue *txsQueue,
) error {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(queue *txsQueue) {
		startReadingChannel(ctx, queue, sentTxsChan)
		wg.Done()
	}(queue)

	wg.Add(1)
	go func(queue *txsQueue) {
		startCheckingTxs(ctx, queue, client)
		wg.Done()
	}(queue)
	wg.Wait()
	return nil
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

func startCheckingTxs(ctx context.Context, queue *txsQueue, client *eth.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if queue.IsEmpty() {
				continue
			}
			tx := queue.First()
			err := waitUntilTxMinedAndCheckForFail(client, tx)
			if err != nil {
				panic(err)
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
