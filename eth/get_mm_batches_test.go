package eth

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetMMBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	testSuiteWithRequestsSending
	client     *TestClient
	commitment *models.MMCommitmentWithTxs
}

func (s *GetMMBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.commitment = &models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				ID: models.CommitmentID{
					BatchID:      models.MakeUint256(1),
					IndexInBatch: 0,
				},
				Type: batchtype.MassMigration,
			},
			CombinedSignature: *mockSignature(s.Assertions),
			Meta: &models.MassMigrationMeta{
				SpokeID:     1,
				TokenID:     models.MakeUint256(0),
				Amount:      models.MakeUint256(400),
				FeeReceiver: 0,
			},
			WithdrawRoot: utils.RandomHash(),
		},
		Transactions: []uint8{0, 0, 0, 0, 32, 4, 0, 0},
	}
}

func (s *GetMMBatchesTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client

	s.StartTxsSending(s.T(), client.TxsChannels.Requests)
}

func (s *GetMMBatchesTestSuite) TearDownTest() {
	s.StopTxsSending()
	s.client.Close()
}

func (s *GetMMBatchesTestSuite) TestGetBatches() {
	batch, err := s.client.SubmitMassMigrationsBatchAndWait(models.NewUint256(1), []models.CommitmentWithTxs{s.commitment})
	s.NoError(err)

	batches, err := s.client.GetBatches(&BatchesFilters{
		FilterByBatchID: func(batchID *models.Uint256) bool {
			return batchID.CmpN(0) > 0
		},
	})
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(batch.ID, batches[0].GetID())
	s.Equal(batchtype.MassMigration, batches[0].GetBase().Type)
}

func (s *GetMMBatchesTestSuite) TestGetTxBatch() {
	batchID := models.MakeUint256(1)
	tx, err := s.client.SubmitMassMigrationsBatch(&batchID, []models.CommitmentWithTxs{s.commitment})
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
	s.Len(decodedMMBatch.Commitments, 1)
	s.EqualValues(event.AccountRoot, decodedMMBatch.AccountTreeRoot)
}

func (s *GetMMBatchesTestSuite) TestGetTxBatch_ReturnsErrorWhenCurrentBatchHasDifferentHash() {
	batchID := models.NewUint256(1)
	tx, err := s.client.SubmitMassMigrationsBatch(batchID, []models.CommitmentWithTxs{s.commitment})
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
