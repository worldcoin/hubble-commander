package eth

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetMMBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	client        *TestClient
	commitments   []models.CommitmentWithTxs
	metas         []models.MassMigrationMeta
	withdrawRoots []common.Hash
}

func (s *GetMMBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.commitments = []models.CommitmentWithTxs{
		{
			TxCommitment: models.TxCommitment{
				CommitmentBase: models.CommitmentBase{
					ID: models.CommitmentID{
						BatchID:      models.MakeUint256(1),
						IndexInBatch: 0,
					},
					Type: batchtype.MassMigration,
				},
				FeeReceiver:       0,
				CombinedSignature: *mockSignature(s.Assertions),
			},
			Transactions: []uint8{0, 0, 0, 0, 32, 4, 0, 0},
		},
	}
	s.metas = []models.MassMigrationMeta{
		{
			SpokeID:     1,
			TokenID:     models.MakeUint256(0),
			Amount:      models.MakeUint256(400),
			FeeReceiver: 0,
		},
	}
	s.withdrawRoots = []common.Hash{
		{1, 2, 3},
	}
}

func (s *GetMMBatchesTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *GetMMBatchesTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *GetMMBatchesTestSuite) TestGetBatches() {
	batch, err := s.client.SubmitMassMigrationsBatchAndWait(models.NewUint256(1), s.commitments, s.metas, s.withdrawRoots)
	s.NoError(err)

	batches, err := s.client.GetBatches(&BatchesFilters{
		FilterByBatchID: func(batchID *models.Uint256) bool {
			return batchID.CmpN(0) > 0
		},
	})
	s.NoError(err)
	s.Len(batches, 1)
	s.EqualValues(batch.ID, batches[0].GetID())
	s.Equal(batch.Type, batches[0].GetBase().Type)
}

func (s *GetMMBatchesTestSuite) TestGetMMBatch_BatchExists() {
	batchID := models.MakeUint256(1)
	tx, err := s.client.SubmitMassMigrationsBatch(&batchID, s.commitments, s.metas, s.withdrawRoots)
	s.NoError(err)
	s.client.GetBackend().Commit()

	transaction, _, err := s.client.Blockchain.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     batchID.ToBig(),
		AccountRoot: getAccountRoot(s.Assertions, s.client),
		BatchType:   uint8(batchtype.MassMigration),
	}

	decodedBatch, err := s.client.getTxBatch(event, transaction, decodeMMCommitments)
	s.NoError(err)
	decodedMMBatch := decodedBatch.ToDecodedTxBatch()
	s.Equal(batchID, decodedMMBatch.ID)
	s.EqualValues(event.BatchType, decodedMMBatch.Type)
	s.Len(decodedMMBatch.Commitments, len(s.commitments))
	s.EqualValues(event.AccountRoot, decodedMMBatch.AccountTreeRoot)
}

func (s *GetMMBatchesTestSuite) TestGetMMBatch_BatchNotExists() {
	tx, err := s.client.SubmitMassMigrationsBatch(models.NewUint256(1), s.commitments, s.metas, s.withdrawRoots)
	s.NoError(err)
	s.client.GetBackend().Commit()

	transaction, _, err := s.client.Blockchain.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     big.NewInt(5),
		AccountRoot: getAccountRoot(s.Assertions, s.client),
		BatchType:   uint8(batchtype.MassMigration),
	}

	batch, err := s.client.getTxBatch(event, transaction, decodeMMCommitments)
	s.Nil(batch)
	s.ErrorIs(err, errBatchAlreadyRolledBack)
}

func (s *GetMMBatchesTestSuite) TestGetMMBatch_DifferentBatchHash() {
	batchID := models.NewUint256(1)
	tx, err := s.client.SubmitMassMigrationsBatch(batchID, s.commitments, s.metas, s.withdrawRoots)
	s.NoError(err)
	s.client.GetBackend().Commit()

	transaction, _, err := s.client.Blockchain.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     batchID.ToBig(),
		AccountRoot: [32]byte{1, 2, 3},
		BatchType:   uint8(batchtype.MassMigration),
	}

	batch, err := s.client.getTxBatch(event, transaction, decodeMMCommitments)
	s.Nil(batch)
	s.ErrorIs(err, errBatchAlreadyRolledBack)
}

func TestGetMMBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetMMBatchesTestSuite))
}
