package eth

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
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
	// TODO filter DepositsFinalised events and return subtreeID as part of DecodedCommitment
	batchEvents, depositEvents, err := c.getBatchEvents(filters)
	if err != nil {
		return nil, err
	}
	logBatchesCount(len(batchEvents))

	depositIndex := 0
	res := make([]DecodedBatch, 0, len(batchEvents))
	for i := range batchEvents {
		if filters.FilterByBatchID != nil && !filters.FilterByBatchID(models.NewUint256FromBig(*batchEvents[i].BatchID)) {
			continue
		}

		tx, _, err := c.Blockchain.GetBackend().TransactionByHash(context.Background(), batchEvents[i].Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitTransfer"].ID) &&
			!bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitCreate2Transfer"].ID) &&
			!bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitDeposits"].ID) {
			continue // TODO handle internal transactions
		}

		var decodedBatch DecodedBatch
		// nolint: exhaustive
		switch batchtype.BatchType(batchEvents[i].BatchType) {
		case batchtype.Transfer, batchtype.Create2Transfer:
			decodedBatch, err = c.getTxBatch(batchEvents[i], tx)
		case batchtype.Deposit:
			decodedBatch, err = c.getDepositBatch(batchEvents[i], depositEvents[depositIndex], tx)
			depositIndex++
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

func (c *Client) getBatchEvents(filters *BatchesFilters) ([]*rollup.RollupNewBatch, []*rollup.RollupDepositsFinalised, error) {
	batchIterator, err := c.Rollup.FilterNewBatch(&bind.FilterOpts{
		Start: filters.StartBlockInclusive,
		End:   filters.EndBlockInclusive,
	})
	if err != nil {
		return nil, nil, err
	}
	events := make([]*rollup.RollupNewBatch, 0)
	for batchIterator.Next() {
		events = append(events, batchIterator.Event)
	}

	depositIterator, err := c.Rollup.FilterDepositsFinalised(&bind.FilterOpts{
		Start: filters.StartBlockInclusive,
		End:   filters.EndBlockInclusive,
	})
	if err != nil {
		return nil, nil, err
	}
	depositEvents := make([]*rollup.RollupDepositsFinalised, 0)
	for depositIterator.Next() {
		depositEvents = append(depositEvents, depositIterator.Event)
	}
	return events, depositEvents, nil
}

func (c *Client) getTxBatch(batchEvent *rollup.RollupNewBatch, tx *types.Transaction) (DecodedBatch, error) {
	batch, err := c.GetBatch(models.NewUint256FromBig(*batchEvent.BatchID))
	if err != nil {
		if err.Error() == MsgInvalidBatchID {
			return nil, errBatchAlreadyRolledBack
		}
		return nil, err
	}

	decodedBatch := &DecodedTxBatch{
		DecodedBatchBase: *NewDecodedBatchBase(batch, tx.Hash(), common.BytesToHash(batchEvent.AccountRoot[:])),
	}
	decodedBatch.Commitments, err = encoder.DecodeBatchCalldata(tx.Data(), &decodedBatch.ID)
	if err != nil {
		return nil, err
	}

	err = decodedBatch.verifyBatchHash()
	if err != nil {
		return nil, err
	}

	err = c.setSubmissionTime(decodedBatch, batchEvent.Raw.BlockNumber)
	if err != nil {
		return nil, err
	}
	return decodedBatch, nil
}

func (c *Client) getDepositBatch(
	batchEvent *rollup.RollupNewBatch,
	depositEvent *rollup.RollupDepositsFinalised,
	tx *types.Transaction,
) (DecodedBatch, error) {
	batch, err := c.GetBatch(models.NewUint256FromBig(*batchEvent.BatchID))
	if err != nil {
		if err.Error() == MsgInvalidBatchID {
			return nil, errBatchAlreadyRolledBack
		}
		return nil, err
	}

	decodedBatch := &DecodedDepositBatch{
		DecodedBatchBase: *NewDecodedBatchBase(batch, tx.Hash(), common.BytesToHash(batchEvent.AccountRoot[:])),
		SubtreeID:        models.MakeUint256FromBig(*depositEvent.SubtreeID),
		PathAtDepth:      uint32(depositEvent.PathToSubTree.Uint64()),
	}
	err = c.setSubmissionTime(decodedBatch, batchEvent.Raw.BlockNumber)
	if err != nil {
		return nil, err
	}
	return decodedBatch, nil
}

func (c *Client) setSubmissionTime(decodedBatch DecodedBatch, blockNumber uint64) error {
	header, err := c.Blockchain.GetBackend().HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return err
	}
	decodedBatch.GetBatch().SubmissionTime = *models.NewTimestamp(time.Unix(int64(header.Time), 0).UTC())
	return nil
}

func logBatchesCount(count int) {
	if count > 0 {
		log.Printf("Found %d batch(es)", count)
	}
}
