package eth

import (
	"bytes"
	"context"
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
	return c.GetBatches(&BatchesFilters{})
}

func (c *Client) GetBatches(filters *BatchesFilters) ([]DecodedBatch, error) {
	batchEvents, depositEvents, err := c.getBatchEvents(filters)
	if err != nil {
		return nil, err
	}
	logBatchesCount(len(batchEvents))

	var rolledBackBatchID *models.Uint256
	depositIndex := 0
	res := make([]DecodedBatch, 0, len(batchEvents))
	for i := range batchEvents {
		event := batchEvents[i]
		if !isAcceptable(models.NewUint256FromBig(*event.BatchID), rolledBackBatchID, filters) {
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
			decodedBatch, err = c.getMMBatch(event, tx)
		case batchtype.Genesis:
			panic("syncing genesis batch should have been skipped")
		}

		if errors.Is(err, errBatchAlreadyRolledBack) {
			rolledBackBatchID = models.NewUint256FromBig(*event.BatchID)
			continue
		}
		if err != nil {
			return nil, err
		}

		if rolledBackBatchID != nil && *rolledBackBatchID == decodedBatch.GetBase().ID {
			rolledBackBatchID = nil
		}
		res = append(res, decodedBatch)
	}

	return res, nil
}

func isAcceptable(batchID, rolledBackBatchID *models.Uint256, filters *BatchesFilters) bool {
	return (filters.FilterByBatchID == nil || filters.FilterByBatchID(batchID)) &&
		(rolledBackBatchID == nil || batchID.Cmp(rolledBackBatchID) <= 0)
}

func (c *Client) getBatchEvents(filters *BatchesFilters) ([]*rollup.RollupNewBatch, []*rollup.RollupDepositsFinalised, error) {
	batchEvents, err := c.getNewBatchEvents(filters)
	if err != nil {
		return nil, nil, err
	}

	depositEvents, err := c.getDepositsFinalisedEvents(filters)
	if err != nil {
		return nil, nil, err
	}

	return batchEvents, depositEvents, nil
}

func (c *Client) getNewBatchEvents(filters *BatchesFilters) ([]*rollup.RollupNewBatch, error) {
	it, err := c.getNewBatchLogIterator(filters)
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	events := make([]*rollup.RollupNewBatch, 0)
	for it.Next() {
		events = append(events, it.Event)
	}

	// Sort for sanity
	sort.Slice(events, func(i, j int) bool {
		return utils.EventBefore(&events[i].Raw, &events[j].Raw)
	})

	return events, nil
}

func (c *Client) getDepositsFinalisedEvents(filters *BatchesFilters) ([]*rollup.RollupDepositsFinalised, error) {
	it, err := c.getDepositsFinalisedLogIterator(filters)
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	events := make([]*rollup.RollupDepositsFinalised, 0)
	for it.Next() {
		events = append(events, it.Event)
	}

	// Sort for sanity
	sort.Slice(events, func(i, j int) bool {
		return utils.EventBefore(&events[i].Raw, &events[j].Raw)
	})

	return events, nil
}

func (c *Client) getNewBatchLogIterator(filters *BatchesFilters) (*rollup.NewBatchIterator, error) {
	it := &rollup.NewBatchIterator{}

	err := c.FilterLogs(c.Rollup.BoundContract, "NewBatch", &bind.FilterOpts{
		Start: filters.StartBlockInclusive,
		End:   filters.EndBlockInclusive,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func (c *Client) getDepositsFinalisedLogIterator(filters *BatchesFilters) (*rollup.DepositsFinalisedIterator, error) {
	it := &rollup.DepositsFinalisedIterator{}

	err := c.FilterLogs(c.Rollup.BoundContract, "DepositsFinalised", &bind.FilterOpts{
		Start: filters.StartBlockInclusive,
		End:   filters.EndBlockInclusive,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func (c *Client) isDirectBatchSubmission(tx *types.Transaction) bool {
	methodID := tx.Data()[:4]
	return bytes.Equal(methodID, c.Rollup.ABI.Methods["submitTransfer"].ID) ||
		bytes.Equal(methodID, c.Rollup.ABI.Methods["submitMassMigration"].ID) ||
		bytes.Equal(methodID, c.Rollup.ABI.Methods["submitCreate2Transfer"].ID) ||
		bytes.Equal(methodID, c.Rollup.ABI.Methods["submitDeposits"].ID)
}

func (c *Client) getTxBatch(batchEvent *rollup.RollupNewBatch, tx *types.Transaction) (DecodedBatch, error) {
	batch, err := c.getBatchDetails(batchEvent)
	if err != nil {
		return nil, err
	}
	commitments, err := encoder.DecodeBatchCalldata(c.Rollup.ABI, tx.Data())
	if err != nil {
		return nil, err
	}
	accountRoot := common.BytesToHash(batchEvent.AccountRoot[:])

	decodedCommitments := decodedTxCommitmentsToCommitments(commitments)
	err = verifyBatchHash(*batch.Hash, accountRoot, decodedCommitments)
	if err != nil {
		return nil, err
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
		Commitments: decodedCommitments,
	}, nil
}

//TODO-sync: replace with function that accepts decodeCalldataFunction
func (c *Client) getMMBatch(batchEvent *rollup.RollupNewBatch, tx *types.Transaction) (DecodedBatch, error) {
	batch, err := c.getBatchDetails(batchEvent)
	if err != nil {
		return nil, err
	}
	decodedCommitments, err := encoder.DecodeMassMigrationBatchCalldata(c.Rollup.ABI, tx.Data())
	if err != nil {
		return nil, err
	}
	accountRoot := common.BytesToHash(batchEvent.AccountRoot[:])

	commitments := decodedMMCommitmentsToCommitments(decodedCommitments)
	err = verifyBatchHash(*batch.Hash, accountRoot, commitments)
	if err != nil {
		return nil, err
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

func verifyBatchHash(batchHash, accountRoot common.Hash, commitments []encoder.GenericCommitment) error {
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
