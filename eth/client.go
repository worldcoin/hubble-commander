package eth

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const TxTimeout = 30 * time.Second

type Client struct {
	account *bind.TransactOpts
	Rollup  *rollup.Rollup
}

func NewClient(rollupContract *rollup.Rollup) *Client {
	return &Client{Rollup: rollupContract}
}

func (c *Client) SubmitTransfer(commitments []*models.Commitment) error {
	count := len(commitments)

	stateRoots := make([][32]byte, 0, count)
	signatures := make([][2]*big.Int, 0, count)
	feeReceivers := make([]*big.Int, 0, count)
	transactions := make([][]byte, 0, count)

	for _, commitment := range commitments {
		stateRoots = append(stateRoots, commitment.PostStateRoot)
		signatures = append(signatures, commitment.CombinedSignature.ToBigIntPointers())
		feeReceivers = append(feeReceivers, &commitment.FeeReceiver.Int)
		transactions = append(transactions, commitment.Transactions)
	}
	_, err := c.Rollup.SubmitTransfer(c.account, stateRoots, signatures, feeReceivers, transactions)
	if err != nil {
		return err
	}

	sink := make(chan *rollup.RollupNewBatch)
	subscription, err := c.Rollup.WatchNewBatch(&bind.WatchOpts{}, sink)
	if err != nil {
		return err
	}
	for {
		select {
		case newBatch := <-sink:
			if c.account.From == newBatch.Committer {
				subscription.Unsubscribe()
				break
			}
		case <-time.After(TxTimeout):
			return fmt.Errorf("timeout")
		}
	}
}
