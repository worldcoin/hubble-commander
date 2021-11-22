package eth

import (
	"bytes"
	"context"
	"github.com/Worldcoin/hubble-commander/metrics"
	"math/big"
	"sort"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const MsgInvalidBatchID = "execution reverted: Batch id greater than total number of batches, invalid batch id"

var errBatchAlreadyRolledBack = errors.New("batch already rolled back")

type BatchesFilters struct {
	StartBlockInclusive uint64
	EndBlockInclusive   *uint64
	FilterByBatchID     func(batchID *models.Uint256) bool
}

func (c *TestClient) GetAllBatches() ([]DecodedBatch, error) {
	return c.GetBatches(&BatchesFilters{}, c.CommanderMetrics)
}

func (c *Client) GetBatches(filters *BatchesFilters, commanderMetrics *metrics.CommanderMetrics) ([]DecodedBatch, error) {
	batchEvents, depositEvents, err := c.getBatchEvents(filters, commanderMetrics)
	if err != nil {
		return nil, err
	}
	logBatchesCount(len(batchEvents))

	depositIndex := 0
	res := make([]DecodedBatch, 0, len(batchEvents))
	for i := range batchEvents {
		event := batchEvents[i]
		if filters.FilterByBatchID != nil && !filters.FilterByBatchID(models.NewUint256FromBig(*event.BatchID)) {
			continue
		}

		tx, _, err := c.Blockchain.GetBackend().TransactionByHash(context.Background(), event.Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !c.isDirectBatchSubmission(tx) {
			continue // TODO handle internal transactions
		}

		var decodedBatch DecodedBatch

		switch batchtype.BatchType(event.BatchType) {
		case batchtype.Transfer, batchtype.Create2Transfer:
			decodedBatch, err = c.getTxBatch(event, tx)
		case batchtype.Deposit:
			decodedBatch, err = c.getDepositBatch(event, depositEvents[depositIndex], tx)
			depositIndex++
		case batchtype.MassMigration:
			panic("syncing MassMigration batches is not supported yet")
		case batchtype.Genesis:
			panic("syncing genesis batch should have been skipped")
		}

		if errors.Is(err, errBatchAlreadyRolledBack) {
			// TODO: handle deposit rollbacks after https://github.com/thehubbleproject/hubble-contracts/issues/671
			continue
		}
		if err != nil {
			return nil, err
		}

		res = append(res, decodedBatch)
	}

	return res, nil
}

func (c *Client) getBatchEvents(filters *BatchesFilters, commanderMetrics *metrics.CommanderMetrics) ([]*rollup.RollupNewBatch, []*rollup.RollupDepositsFinalised, error) {
	newBatchIterator, err := c.getNewBatchLogIterator(filters, commanderMetrics)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = newBatchIterator.Close() }()

	batchEvents := make([]*rollup.RollupNewBatch, 0)
	for newBatchIterator.Next() {
		batchEvents = append(batchEvents, newBatchIterator.Event)
	}

	// Sort for sanity
	sort.Slice(batchEvents, func(i, j int) bool {
		return utils.EventBefore(&batchEvents[i].Raw, &batchEvents[j].Raw)
	})

	depositsFinalisedIterator, err := c.getDepositsFinalisedLogIterator(filters, commanderMetrics)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = depositsFinalisedIterator.Close() }()

	depositEvents := make([]*rollup.RollupDepositsFinalised, 0)
	for depositsFinalisedIterator.Next() {
		depositEvents = append(depositEvents, depositsFinalisedIterator.Event)
	}

	// Sort for sanity
	sort.Slice(depositEvents, func(i, j int) bool {
		return utils.EventBefore(&batchEvents[i].Raw, &batchEvents[j].Raw)
	})

	return batchEvents, depositEvents, nil
}

func (c *Client) getNewBatchLogIterator(filters *BatchesFilters, commanderMetrics *metrics.CommanderMetrics) (*rollup.RollupNewBatchIterator, error) {
	var it *rollup.RollupNewBatchIterator

	duration, err := metrics.MeasureDuration(func() (err error) {
		it, err = c.Rollup.FilterNewBatch(&bind.FilterOpts{
			Start: filters.StartBlockInclusive,
			End:   filters.EndBlockInclusive,
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	commanderMetrics.SaveBlockchainCallDurationMeasurement(*duration, metrics.NewBatchLogRetrievalCall)

	return it, nil
}
func (c *Client) getDepositsFinalisedLogIterator(filters *BatchesFilters, commanderMetrics *metrics.CommanderMetrics) (*rollup.RollupDepositsFinalisedIterator, error) {
	var it *rollup.RollupDepositsFinalisedIterator

	duration, err := metrics.MeasureDuration(func() (err error) {
		it, err = c.Rollup.FilterDepositsFinalised(&bind.FilterOpts{
			Start: filters.StartBlockInclusive,
			End:   filters.EndBlockInclusive,
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	commanderMetrics.SaveBlockchainCallDurationMeasurement(*duration, metrics.DepositsFinalisedLogRetrievalCall)

	return it, nil
}

func (c *Client) isDirectBatchSubmission(tx *types.Transaction) bool {
	methodID := tx.Data()[:4]
	return bytes.Equal(methodID, c.RollupABI.Methods["submitTransfer"].ID) ||
		bytes.Equal(methodID, c.RollupABI.Methods["submitCreate2Transfer"].ID) ||
		bytes.Equal(methodID, c.RollupABI.Methods["submitDeposits"].ID)
}

func (c *Client) getTxBatch(batchEvent *rollup.RollupNewBatch, tx *types.Transaction) (DecodedBatch, error) {
	batch, err := c.getBatchDetails(batchEvent)
	if err != nil {
		return nil, err
	}
	commitments, err := encoder.DecodeBatchCalldata(tx.Data(), &batch.ID)
	if err != nil {
		return nil, err
	}
	accountRoot := common.BytesToHash(batchEvent.AccountRoot[:])

	if vErr := verifyBatchHash(*batch.Hash, accountRoot, commitments); vErr != nil {
		return nil, vErr
	}

	timestamp, err := c.getBlockTimestamp(batchEvent.Raw.BlockNumber)
	if err != nil {
		return nil, err
	}

	return &DecodedTxBatch{
		DecodedBatchBase: *NewDecodedBatchBase(
			batch,
			tx.Hash(),
			accountRoot,
			timestamp,
		),
		Commitments: commitments,
	}, nil
}

func (c *Client) getDepositBatch(
	batchEvent *rollup.RollupNewBatch,
	depositEvent *rollup.RollupDepositsFinalised,
	tx *types.Transaction,
) (DecodedBatch, error) {
	batch, err := c.getBatchDetails(batchEvent)
	if err != nil {
		return nil, err
	}
	accountRoot := common.BytesToHash(batchEvent.AccountRoot[:])

	timestamp, err := c.getBlockTimestamp(batchEvent.Raw.BlockNumber)
	if err != nil {
		return nil, err
	}

	return &DecodedDepositBatch{
		DecodedBatchBase: *NewDecodedBatchBase(batch, tx.Hash(), accountRoot, timestamp),
		SubtreeID:        models.MakeUint256FromBig(*depositEvent.SubtreeID),
		PathAtDepth:      uint32(depositEvent.PathToSubTree.Uint64()),
	}, nil
}

func (c *Client) getBatchDetails(batchEvent *rollup.RollupNewBatch) (*models.Batch, error) {
	batchID := models.NewUint256FromBig(*batchEvent.BatchID)
	batch, err := c.GetBatch(batchID)
	if err != nil && err.Error() == MsgInvalidBatchID {
		return nil, errBatchAlreadyRolledBack
	}
	return batch, err
}

func (c *Client) getBlockTimestamp(blockNumber uint64) (*models.Timestamp, error) {
	header, err := c.Blockchain.GetBackend().HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, err
	}
	utcTime := time.Unix(int64(header.Time), 0).UTC()
	return models.NewTimestamp(utcTime), nil
}

func verifyBatchHash(batchHash, accountRoot common.Hash, commitments []encoder.DecodedCommitment) error {
	leafHashes := make([]common.Hash, 0, len(commitments))
	for i := range commitments {
		leafHashes = append(leafHashes, commitments[i].LeafHash(accountRoot))
	}
	tree, err := merkletree.NewMerkleTree(leafHashes)
	if err != nil {
		return err
	}

	if tree.Root() != batchHash {
		return errBatchAlreadyRolledBack
	}
	return nil
}

func logBatchesCount(count int) {
	if count > 0 {
		log.Printf("Found %d batch(es)", count)
	}
}
