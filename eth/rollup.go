package eth

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) rollup() *RollupSessionBuilder {
	return &RollupSessionBuilder{rollup.RollupSession{
		Contract:     c.Rollup,
		TransactOpts: c.account,
	}}
}

func (c *Client) SubmitTransfersBatch(commitments []models.Commitment) (batch *models.Batch, accountTreeRoot *common.Hash, err error) {
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
				return c.handleNewBatchEvent(newBatch)
			}
		case <-time.After(*c.config.txTimeout):
			return nil, nil, fmt.Errorf("timeout")
		}
	}
}

func (c *Client) GetBatch(batchID *models.Uint256) (*models.Batch, error) {
	batch, err := c.Rollup.GetBatch(nil, &batchID.Int)
	if err != nil {
		return nil, err
	}
	meta := encoder.DecodeMeta(batch.Meta)
	return &models.Batch{
		Hash:              common.BytesToHash(batch.CommitmentRoot[:]),
		ID:                *batchID,
		FinalisationBlock: meta.FinaliseOn,
	}, nil
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

	for _, commitment := range commitments {
		stateRoots = append(stateRoots, commitment.PostStateRoot)
		signatures = append(signatures, commitment.CombinedSignature.ToBigIntPointers())
		feeReceivers = append(feeReceivers, new(big.Int).SetUint64(uint64(commitment.FeeReceiver)))
		transactions = append(transactions, commitment.Transactions)
	}
	return
}
