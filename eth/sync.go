package eth

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) GetBatches() ([]DecodedBatch, error) {
	it, err := c.Rollup.FilterNewBatch(nil)
	if err != nil {
		return nil, err
	}

	res := make([]DecodedBatch, 0)
	for it.Next() {
		address := it.Event.Raw.Address
		txHash := it.Event.Raw.TxHash

		// TODO: handle internal transactions
		tx, _, err := c.ChainConnection.GetBackend().TransactionByHash(context.Background(), txHash)
		if err != nil {
			return nil, err
		}

		if *tx.To() != address {
			return nil, fmt.Errorf("log address is different from the contract address")
		}

		commitments, err := encoder.DecodeBatch(tx.Data()[4:])
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
			commitments:       commitments,
		})
	}

	return res, nil
}

type DecodedBatch struct {
	Hash              common.Hash
	Type              txtype.TransactionType
	ID                models.Uint256
	FinalisationBlock uint32 // nolint:misspell
	commitments       []encoder.DecodedCommitment
}
