package eth

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SubmitBatchFunc func(commitments []models.Commitment) (*types.Transaction, error)

func (c *Client) SubmitTransfersBatch(commitments []models.Commitment) (
	batch *models.Batch,
	accountTreeRoot *common.Hash,
	err error,
) {
	return c.submitBatch(commitments, func(commitments []models.Commitment) (*types.Transaction, error) {
		return c.rollup().
			WithValue(c.config.stakeAmount.Int).
			SubmitTransfer(parseCommitments(commitments))
	})
}

func (c *Client) SubmitCreate2TransfersBatch(commitments []models.Commitment) (
	batch *models.Batch,
	accountTreeRoot *common.Hash,
	err error,
) {
	return c.submitBatch(commitments, func(commitments []models.Commitment) (*types.Transaction, error) {
		return c.rollup().
			WithValue(c.config.stakeAmount.Int).
			SubmitCreate2Transfer(parseCommitments(commitments))
	})
}

func (c *Client) submitBatch(
	commitments []models.Commitment,
	submit SubmitBatchFunc,
) (batch *models.Batch, accountTreeRoot *common.Hash, err error) {
	sink := make(chan *rollup.RollupNewBatch)
	subscription, err := c.Rollup.WatchNewBatch(&bind.WatchOpts{}, sink)
	if err != nil {
		return
	}
	defer subscription.Unsubscribe()

	tx, err := submit(commitments)
	if err != nil {
		return
	}

	for {
		select {
		case newBatch := <-sink:
			if newBatch.Raw.TxHash == tx.Hash() {
				return c.handleNewBatchEvent(newBatch)
			}
		case <-time.After(*c.config.txTimeout):
			return nil, nil, fmt.Errorf("timeout")
		}
	}
}

func (c *Client) handleNewBatchEvent(event *rollup.RollupNewBatch) (*models.Batch, *common.Hash, error) {
	batch, err := c.GetBatch(models.NewUint256FromBig(*event.BatchID))
	if err != nil {
		return nil, nil, err
	}
	accountRoot := common.BytesToHash(event.AccountRoot[:])
	return batch, &accountRoot, nil
}

func parseCommitments(commitments []models.Commitment) (
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

	for i := range commitments {
		stateRoots = append(stateRoots, commitments[i].PostStateRoot)
		signatures = append(signatures, commitments[i].CombinedSignature.ToBigIntPointers())
		feeReceivers = append(feeReceivers, new(big.Int).SetUint64(uint64(commitments[i].FeeReceiver)))
		transactions = append(transactions, commitments[i].Transactions)
	}
	return
}
