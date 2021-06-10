package eth

import (
	"bytes"
	"context"
	"strings"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (c *Client) GetBatches(opts *bind.FilterOpts) ([]DecodedBatch, error) {
	it, err := c.Rollup.FilterNewBatch(opts)
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

		if !bytes.Equal(tx.Data()[:4], rollupAbi.Methods["submitTransfer"].ID) &&
			!bytes.Equal(tx.Data()[:4], rollupAbi.Methods["submitCreate2Transfer"].ID) {
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

		batch.TransactionHash = txHash
		batch.AccountTreeRoot = &accountRoot

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
