package eth

import (
	"bytes"
	"context"
	"strings"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (c *Client) GetBatches() ([]DecodedBatch, error) {
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

		commitments, err := encoder.DecodeBatch(tx.Data())
		if err != nil {
			return nil, err
		}

		meta, err := c.Rollup.Batches(nil, it.Event.BatchID)
		if err != nil {
			return nil, err
		}

		res = append(res, DecodedBatch{
			Hash:              common.BytesToHash(meta.CommitmentRoot[:]),
			Type:              txtype.TransactionType(it.Event.BatchType),
			ID:                models.MakeUint256FromBig(*it.Event.BatchID),
			FinalisationBlock: encoder.DecodeMeta(meta.Meta).FinaliseOn,
			AccountRoot:       common.BytesToHash(it.Event.AccountRoot[:]),
			Commitments:       commitments,
		})
	}

	return res, nil
}

type DecodedBatch struct {
	Hash              common.Hash
	Type              txtype.TransactionType
	ID                models.Uint256
	FinalisationBlock uint32
	AccountRoot       common.Hash
	Commitments       []encoder.DecodedCommitment
}
