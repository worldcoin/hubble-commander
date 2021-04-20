package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) rollup() *RollupSessionBuilder {
	return &RollupSessionBuilder{rollup.RollupSession{
		Contract:     c.Rollup,
		TransactOpts: *c.ChainConnection.GetAccount(),
	}}
}

func (c *Client) GetBatch(batchID *models.Uint256) (*models.Batch, error) {
	batch, err := c.Rollup.GetBatch(nil, &batchID.Int)
	if err != nil {
		return nil, err
	}
	meta := encoder.DecodeMeta(batch.Meta)
	return &models.Batch{
		Hash:              common.BytesToHash(batch.CommitmentRoot[:]),
		Type:              meta.BatchType,
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

	for i := range commitments {
		stateRoots = append(stateRoots, commitments[i].PostStateRoot)
		signatures = append(signatures, commitments[i].CombinedSignature.ToBigIntPointers())
		feeReceivers = append(feeReceivers, new(big.Int).SetUint64(uint64(commitments[i].FeeReceiver)))
		transactions = append(transactions, commitments[i].Transactions)
	}
	return
}
