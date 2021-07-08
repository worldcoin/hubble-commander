package eth

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) GetBatches(opts *bind.FilterOpts) ([]DecodedBatch, error) {
	it, err := c.Rollup.FilterNewBatch(opts)
	if err != nil {
		return nil, err
	}

	res := make([]DecodedBatch, 0)
	for it.Next() {
		txHash := it.Event.Raw.TxHash

		tx, _, err := c.ChainConnection.GetBackend().TransactionByHash(context.Background(), txHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitTransfer"].ID) &&
			!bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitCreate2Transfer"].ID) {
			continue // TODO handle internal transactions
		}

		commitments, err := encoder.DecodeBatchCalldata(tx.Data())
		if err != nil {
			return nil, err
		}

		batch, err := c.GetBatch(models.NewUint256FromBig(*it.Event.BatchID))
		if err != nil {
			return nil, err
		}

		accountRoot := common.BytesToHash(it.Event.AccountRoot[:])

		header, err := c.ChainConnection.GetBackend().HeaderByNumber(context.Background(), new(big.Int).SetUint64(it.Event.Raw.BlockNumber))
		if err != nil {
			return nil, err
		}

		batch.TransactionHash = txHash
		batch.AccountTreeRoot = &accountRoot
		batch.SubmissionTime = ref.Time(time.Unix(int64(header.Time), 0).UTC())

		res = append(res, DecodedBatch{
			Batch:       *batch,
			Commitments: commitments,
		})
	}

	return res, nil
}

type DecodedBatch struct {
	models.Batch
	Commitments []encoder.DecodedCommitment
}
