package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxsTracker struct {
	TxsChan chan *types.Transaction
	*TxsSender
}

func NewTxTracker(ethClient *eth.Client, txsChan chan *types.Transaction) *TxsTracker {
	return &TxsTracker{
		TxsSender: newTxRequestsSender(ethClient),
		TxsChan:   txsChan,
	}
}

func (t *TxsTracker) StartTracking(ctx context.Context) error {
	sentTxs := make([]*types.Transaction, 0)

	for {
		select {
		case <-ctx.Done():
			return nil
		case reqParams := <-t.requestsChan:
			_, err := t.sendTransaction(reqParams)
			if err != nil {
				panic(err)
			}
			continue
		case tx := <-t.TxsChan:
			sentTxs = append(sentTxs, tx)
			continue
		default:
			break
		}

		if len(sentTxs) == 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		isMined, receipt, err := t.client.IsTxMined(sentTxs[0])
		if err != nil {
			panic(err)
		}
		if !isMined {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		err = t.checkTxForFail(sentTxs[0], receipt)
		if err != nil {
			panic(err)
		}
		sentTxs = sentTxs[1:]
	}
	return nil
}

func (t *TxsTracker) checkTxForFail(tx *types.Transaction, receipt *types.Receipt) error {
	if receipt.Status == 1 {
		return nil
	}
	err := t.client.GetRevertMessage(tx, receipt)
	return fmt.Errorf("%w tx_hash=%s", err, tx.Hash().String())
}
