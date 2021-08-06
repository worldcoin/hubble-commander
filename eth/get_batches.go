package eth

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const MsgInvalidBatchID = "execution reverted: Batch id greater than total number of batches, invalid batch id"

var errBatchAlreadyRolledBack = errors.New("batch already rolled back")

func (c *Client) GetBatches(opts *bind.FilterOpts) ([]DecodedBatch, error) {
	return c.GetBatchesInRange(opts, nil, nil)
}

func (c *Client) GetBatchesInRange(opts *bind.FilterOpts, startID, endID *models.Uint256) ([]DecodedBatch, error) {
	it, err := c.Rollup.FilterNewBatch(opts)
	if err != nil {
		return nil, err
	}

	// TODO Convert iterator to slice and print "Found X batches" log here

	res := make([]DecodedBatch, 0)
	for it.Next() {
		batchID := models.NewUint256FromBig(*it.Event.BatchID)
		if startID != nil && batchID.Cmp(startID) <= 0 {
			log.Printf("Batch #%d already synced. Skipping...", batchID.Uint64())
			continue
		}
		if endID != nil && batchID.Cmp(endID) >= 0 {
			log.Printf("Batch #%d after dispute. Skipping...", batchID.Uint64())
			continue
		}

		txHash := it.Event.Raw.TxHash

		tx, _, err := c.ChainConnection.GetBackend().TransactionByHash(context.Background(), txHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitTransfer"].ID) &&
			!bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitCreate2Transfer"].ID) {
			continue // TODO handle internal transactions
		}

		batch, err := c.getBatchIfExists(it.Event, tx)
		if errors.Is(err, errBatchAlreadyRolledBack) {
			continue
		}
		if err != nil {
			return nil, err
		}

		header, err := c.ChainConnection.GetBackend().HeaderByNumber(context.Background(), new(big.Int).SetUint64(it.Event.Raw.BlockNumber))
		if err != nil {
			return nil, err
		}

		batch.TransactionHash = txHash
		batch.SubmissionTime = models.NewTimestamp(time.Unix(int64(header.Time), 0).UTC())

		res = append(res, *batch)
	}

	return res, nil
}

func (c *Client) getBatchIfExists(event *rollup.RollupNewBatch, tx *types.Transaction) (*DecodedBatch, error) {
	batch, err := c.GetBatch(models.NewUint256FromBig(*event.BatchID))
	if err != nil {
		if err.Error() == MsgInvalidBatchID {
			return nil, errBatchAlreadyRolledBack
		}
		return nil, err
	}

	accountRoot := common.BytesToHash(event.AccountRoot[:])
	batch.AccountTreeRoot = &accountRoot

	commitments, err := encoder.DecodeBatchCalldata(tx.Data())
	if err != nil {
		return nil, err
	}

	err = verifyBatchHash(batch, commitments)
	if err != nil {
		return nil, err
	}

	return &DecodedBatch{
		Batch:       *batch,
		Commitments: commitments,
	}, nil
}

func verifyBatchHash(batch *models.Batch, commitments []encoder.DecodedCommitment) error {
	leafHashes := make([]common.Hash, 0, len(commitments))
	for i := range commitments {
		leafHashes = append(leafHashes, commitments[i].LeafHash(*batch.AccountTreeRoot))
	}
	tree, err := merkletree.NewMerkleTree(leafHashes)
	if err != nil {
		return err
	}

	if tree.Root() != *batch.Hash {
		return errBatchAlreadyRolledBack
	}
	return nil
}

type DecodedBatch struct {
	models.Batch
	Commitments []encoder.DecodedCommitment
}
