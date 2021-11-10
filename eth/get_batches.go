package eth

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
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
	batches, err := c.GetBatches(&BatchesFilters{})
	if err != nil {
		return nil, err
	}

	return batches, nil
}

func (c *Client) GetBatches(filters *BatchesFilters) ([]DecodedBatch, error) {
	it, err := c.Rollup.FilterNewBatch(&bind.FilterOpts{
		Start: filters.StartBlockInclusive,
		End:   filters.EndBlockInclusive,
	})
	if err != nil {
		return nil, err
	}

	events := make([]*rollup.RollupNewBatch, 0)
	for it.Next() {
		events = append(events, it.Event)
	}
	logBatchesCount(len(events))

	res := make([]DecodedBatch, 0, len(events))
	for i := range events {
		if filters.FilterByBatchID != nil && !filters.FilterByBatchID(models.NewUint256FromBig(*events[i].BatchID)) {
			continue
		}

		txHash := events[i].Raw.TxHash
		tx, _, err := c.Blockchain.GetBackend().TransactionByHash(context.Background(), txHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitTransfer"].ID) &&
			!bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitCreate2Transfer"].ID) &&
			!bytes.Equal(tx.Data()[:4], c.RollupABI.Methods["submitDeposits"].ID) {
			continue // TODO handle internal transactions
		}

		decodedBatch, err := c.getBatchIfExists(events[i], tx)
		if errors.Is(err, errBatchAlreadyRolledBack) {
			// TODO: handle deposit rollbacks after https://github.com/thehubbleproject/hubble-contracts/issues/671
			continue
		}
		if err != nil {
			return nil, err
		}

		header, err := c.Blockchain.GetBackend().HeaderByNumber(context.Background(), new(big.Int).SetUint64(events[i].Raw.BlockNumber))
		if err != nil {
			return nil, err
		}

		decodedBatch.GetBatch().TransactionHash = txHash
		decodedBatch.GetBatch().SubmissionTime = models.NewTimestamp(time.Unix(int64(header.Time), 0).UTC())

		res = append(res, decodedBatch)
	}

	return res, nil
}

func (c *Client) getBatchIfExists(event *rollup.RollupNewBatch, tx *types.Transaction) (DecodedBatch, error) {
	batch, err := c.GetBatch(models.NewUint256FromBig(*event.BatchID))
	if err != nil {
		if err.Error() == MsgInvalidBatchID {
			return nil, errBatchAlreadyRolledBack
		}
		return nil, err
	}

	accountRoot := common.BytesToHash(event.AccountRoot[:])
	batch.AccountTreeRoot = &accountRoot

	decodedBatch := newDecodedBatch(batch)
	err = decodedBatch.SetCalldata(tx.Data())
	if err != nil {
		return nil, err
	}

	err = decodedBatch.verifyBatchHash()
	if err != nil {
		return nil, err
	}

	return decodedBatch, nil
}

func logBatchesCount(count int) {
	if count > 0 {
		log.Printf("Found %d batch(es)", count)
	}
}
