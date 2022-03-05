package tracker

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (t *Tracker) TrackTxs(ctx context.Context) error {
	wg := sync.WaitGroup{}
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 3)
	wg.Add(1)
	go func() {
		t.startReadingTxsChanLoop()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errChan <- t.startCheckingTxsLoop(subCtx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errChan <- t.sendRequestedTxs(subCtx)
		wg.Done()
	}()

	err := <-errChan
	cancel()
	wg.Wait()
	return err
}

func (t *Tracker) startReadingTxsChanLoop() {
	for {
		tx, isOpen := <-t.txsChan
		if !isOpen {
			return
		}
		t.addTx(tx)
	}
}

func (t *Tracker) startCheckingTxsLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if t.isEmptyTxsQueue() {
				continue
			}
			tx := t.firstTx()
			err := t.waitUntilTxMinedAndCheckForFail(tx)
			if err != nil {
				return err
			}

			t.removeFirstTx()
		}
	}
}

func (t *Tracker) waitUntilTxMinedAndCheckForFail(tx *types.Transaction) error {
	receipt, err := t.client.WaitToBeMined(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if receipt.Status == 1 {
		return nil
	}
	err = t.client.GetRevertMessage(tx, receipt)
	return fmt.Errorf("%w txHash=%s", err, tx.Hash().String())
}
