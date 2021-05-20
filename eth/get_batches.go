package eth

import (
	"bytes"
	"context"
	"strings"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (c *Client) GetBatches(latestBatchSumissionBlock uint32) ([]DecodedBatch, error) {
	it, err := c.Rollup.FilterNewBatch(nil)
	if err != nil {
		return nil, err
	}

	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]DecodedBatch, 0)
	for it.Next() {
		txHash := it.Event.Raw.TxHash

		tx, _, err := c.ChainConnection.GetBackend().TransactionByHash(context.Background(), txHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], rollupAbi.Methods["submitTransfer"].ID) {
			continue // TODO handle internal transactions
		}

		txReceipt, err := c.ChainConnection.GetBackend().TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return nil, err
		}

		if txReceipt.BlockNumber.Uint64() >= uint64(latestBatchSumissionBlock) {
			commitments, err := encoder.DecodeBatchCalldata(tx.Data())
			if err != nil {
				return nil, err
			}

			batch, err := c.GetBatch(models.NewUint256FromBig(*it.Event.BatchID))
			if err != nil {
				return nil, err
			}

			res = append(res, DecodedBatch{
				Batch:       *batch,
				AccountRoot: common.BytesToHash(it.Event.AccountRoot[:]),
				Commitments: commitments,
			})
		}
	}

	return res, nil
}

type DecodedBatch struct {
	models.Batch
	AccountRoot common.Hash
	Commitments []encoder.DecodedCommitment
}
