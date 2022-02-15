package tracker

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func StartTrackingSentTxs(ctx context.Context, client *eth.Client, txsHashChan <-chan *types.Transaction) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case tx := <-txsHashChan:
			err := waitUntilTxMinedAndCheckForFail(client, tx)
			if err != nil {
				panic(err)
			}
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
