package eth

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	client      *TestClient
	commitments []models.CommitmentWithTxs
}

func (s *GetBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.commitments = []models.CommitmentWithTxs{
		{
			TxCommitment: models.TxCommitment{
				CommitmentBase: models.CommitmentBase{
					ID: models.CommitmentID{
						BatchID:      models.MakeUint256(1),
						IndexInBatch: 0,
					},
					Type: batchtype.Transfer,
				},
				FeeReceiver:       0,
				CombinedSignature: *s.mockSignature(),
			},
			Transactions: []uint8{0, 0, 0, 0, 0, 0, 0, 1, 32, 4, 0, 0},
		},
		{
			TxCommitment: models.TxCommitment{
				CommitmentBase: models.CommitmentBase{
					ID: models.CommitmentID{
						BatchID:      models.MakeUint256(2),
						IndexInBatch: 0,
					},
					Type: batchtype.Transfer,
				},
				FeeReceiver:       0,
				CombinedSignature: *s.mockSignature(),
			},
			Transactions: []uint8{0, 0, 1, 0, 0, 0, 0, 0, 32, 1, 0, 0},
		},
	}
}

func (s *GetBatchesTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *GetBatchesTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *GetBatchesTestSuite) TestGetAllBatches() {
	batch1, err := s.client.SubmitTransfersBatchAndWait([]models.CommitmentWithTxs{s.commitments[0]})
	s.NoError(err)
	batch2, err := s.client.SubmitTransfersBatchAndWait([]models.CommitmentWithTxs{s.commitments[1]})
	s.NoError(err)

	batches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(batches, 2)
	s.Equal(batch1.ID, batches[0].ID)
	s.Equal(batch2.ID, batches[1].ID)
}

func (s *GetBatchesTestSuite) TestGetBatches_FiltersByBlockNumber() {
	finalisationBlocks, err := s.client.GetBlocksToFinalise()
	s.NoError(err)

	batch1, err := s.client.SubmitTransfersBatchAndWait([]models.CommitmentWithTxs{s.commitments[0]})
	s.NoError(err)
	batch2, err := s.client.SubmitTransfersBatchAndWait([]models.CommitmentWithTxs{s.commitments[1]})
	s.NoError(err)

	batches, err := s.client.GetBatches(&BatchesFilters{
		StartBlockInclusive: uint64(*batch1.FinalisationBlock - uint32(*finalisationBlocks) + 1),
	})
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(batch2.ID, batches[0].GetBatch().ID)
	s.NotEqual(common.Hash{}, batches[0].GetBatch().TransactionHash)
	s.Equal(s.getAccountRoot(), *batches[0].GetBatch().AccountTreeRoot)
}

func (s *GetBatchesTestSuite) TestGetBatches_FiltersByBatchID() {
	batch1, err := s.client.SubmitTransfersBatchAndWait([]models.CommitmentWithTxs{s.commitments[0]})
	s.NoError(err)
	batch2, err := s.client.SubmitTransfersBatchAndWait([]models.CommitmentWithTxs{s.commitments[1]})
	s.NoError(err)

	batches, err := s.client.GetBatches(&BatchesFilters{
		FilterByBatchID: func(batchID *models.Uint256) bool {
			return batchID.CmpN(0) > 0 && batchID.Cmp(&batch2.ID) < 0
		},
	})
	s.NoError(err)
	s.Len(batches, 1)
	s.EqualValues(batch1.ID, batches[0].GetBatch().ID)
}

func (s *GetBatchesTestSuite) TestGetBatchIfExists_BatchExists() {
	tx, err := s.client.SubmitTransfersBatch(s.commitments)
	s.NoError(err)
	s.client.GetBackend().Commit()

	transaction, _, err := s.client.Blockchain.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     big.NewInt(1),
		AccountRoot: s.getAccountRoot(),
		BatchType:   1,
	}

	decodedBatch, err := s.client.getBatchIfExists(event, transaction)
	s.NoError(err)
	batch := decodedBatch.ToDecodedTxBatch()
	s.Equal(models.MakeUint256(1), batch.ID)
	s.Len(batch.Commitments, len(s.commitments))
	s.EqualValues(event.AccountRoot, *batch.AccountTreeRoot)
}

func (s *GetBatchesTestSuite) TestGetBatchIfExists_BatchNotExists() {
	tx, err := s.client.SubmitTransfersBatch(s.commitments)
	s.NoError(err)
	s.client.GetBackend().Commit()

	transaction, _, err := s.client.Blockchain.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     big.NewInt(5),
		AccountRoot: s.getAccountRoot(),
		BatchType:   1,
	}

	batch, err := s.client.getBatchIfExists(event, transaction)
	s.Nil(batch)
	s.ErrorIs(err, errBatchAlreadyRolledBack)
}

func (s *GetBatchesTestSuite) TestGetBatchIfExists_DifferentBatchHash() {
	tx, err := s.client.SubmitTransfersBatch(s.commitments)
	s.NoError(err)
	s.client.GetBackend().Commit()

	transaction, _, err := s.client.Blockchain.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     big.NewInt(1),
		AccountRoot: [32]byte{1, 2, 3},
		BatchType:   1,
	}

	batch, err := s.client.getBatchIfExists(event, transaction)
	s.Nil(batch)
	s.ErrorIs(err, errBatchAlreadyRolledBack)
}

func (s *GetBatchesTestSuite) getAccountRoot() common.Hash {
	rawAccountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	return common.BytesToHash(rawAccountRoot[:])
}

func (s *GetBatchesTestSuite) mockSignature() *models.Signature {
	wallet, err := bls.NewRandomWallet(bls.Domain{1, 2, 3, 4})
	s.NoError(err)
	signature, err := wallet.Sign(utils.RandomBytes(4))
	s.NoError(err)
	return signature.ModelsSignature()
}

func TestGetBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchesTestSuite))
}
