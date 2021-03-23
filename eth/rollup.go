package eth

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Client) rollup() *RollupSessionBuilder {
	return &RollupSessionBuilder{rollup.RollupSession{
		Contract:     c.Rollup,
		TransactOpts: c.account,
	}}
}

func (c *Client) SubmitTransfersBatch(commitments []*models.Commitment) (batchID *models.Uint256, err error) {
	sink := make(chan *rollup.RollupNewBatch)
	subscription, err := c.Rollup.WatchNewBatch(&bind.WatchOpts{}, sink)
	if err != nil {
		return
	}
	defer subscription.Unsubscribe()

	tx, err := c.rollup().
		WithValue(c.config.stakeAmount.Int).
		SubmitTransfer(parseCommitments(commitments))
	if err != nil {
		return
	}

	for {
		select {
		case newBatch := <-sink:
			if newBatch.Raw.TxHash == tx.Hash() {
				return models.NewUint256FromBig(*newBatch.BatchID), nil
			}
		case <-time.After(*c.config.txTimeout):
			return nil, fmt.Errorf("timeout")
		}
	}
}

func parseCommitments(commitments []*models.Commitment) (
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	feeReceivers []*big.Int,
	transactions [][]byte,
) {
	count := len(commitments)

	stateRoots = make([][32]byte, 0, count)
	signatures = make([][2]*big.Int, 0, count)
	feeReceivers = make([]*big.Int, 0, count)
	transactions = make([][]byte, 0, count)

	for _, commitment := range commitments {
		stateRoots = append(stateRoots, commitment.PostStateRoot)
		signatures = append(signatures, commitment.CombinedSignature.ToBigIntPointers())
		feeReceivers = append(feeReceivers, new(big.Int).SetUint64(uint64(commitment.FeeReceiver)))
		transactions = append(transactions, commitment.Transactions)
	}
	return
}
