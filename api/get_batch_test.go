package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	api                 *API
	storage             *st.TestStorage
	testClient          *eth.TestClient
	txCommitment        *models.TxCommitment
	mmCommitment        *models.MMCommitment
	depositCommitment   *models.DepositCommitment
	batch               *models.Batch
	batchNotFoundAPIErr *APIError
}

func (s *GetBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage, client: s.testClient.Client}

	s.batch = &models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(42000),
		AccountTreeRoot:   utils.NewRandomHash(),
		SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
	}

	s.txCommitment = &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      s.batch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       1,
		CombinedSignature: models.MakeRandomSignature(),
		BodyHash:          utils.NewRandomHash(),
	}

	s.mmCommitment = &models.MMCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      s.batch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.MassMigration,
			PostStateRoot: utils.RandomHash(),
		},
		CombinedSignature: models.MakeRandomSignature(),
		BodyHash:          utils.NewRandomHash(),
		Meta: &models.MassMigrationMeta{
			SpokeID:     1,
			TokenID:     models.MakeUint256(1),
			Amount:      models.MakeUint256(1),
			FeeReceiver: 1,
		},
		WithdrawRoot: utils.RandomHash(),
	}

	s.depositCommitment = &models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      s.batch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.Deposit,
			PostStateRoot: utils.RandomHash(),
		},
		SubtreeID:   models.MakeUint256(1),
		SubtreeRoot: utils.RandomHash(),
		Deposits: []models.PendingDeposit{
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(1),
					DepositIndex: models.MakeUint256(0),
				},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(1),
				L2Amount:   models.MakeUint256(1000),
			},
		},
	}

	s.batchNotFoundAPIErr = &APIError{
		Code:    30000,
		Message: "batch not found",
	}
}

func (s *GetBatchTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_TxBatch() {
	s.addStateLeaf()
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.txCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result, s.batch.Type)
	s.validateTxCommitment(result, s.txCommitment)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_MassMigrationBatch() {
	s.addStateLeaf()

	s.batch.Type = batchtype.MassMigration
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.mmCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result, s.batch.Type)
	s.validateMMCommitment(result, s.mmCommitment)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_DepositBatch() {
	s.addStateLeaf()

	s.batch.Type = batchtype.Deposit
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.depositCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result, s.batch.Type)
	s.validateDepositCommitment(result)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_GenesisBatch() {
	genesisBatch := models.Batch{
		ID:                models.MakeUint256(0),
		Type:              batchtype.Genesis,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(10),
	}
	err := s.storage.AddBatch(&genesisBatch)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*genesisBatch.Hash)
	s.NoError(err)
	s.Equal(genesisBatch.Hash, result.Hash)
	s.Equal(genesisBatch.Type, result.Type)
	s.Equal(*genesisBatch.FinalisationBlock, *result.SubmissionBlock)
	s.Nil(result.Commitments)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NonexistentBatch() {
	result, err := s.api.GetBatchByHash(utils.RandomHash())
	s.Equal(s.batchNotFoundAPIErr, err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID_TxBatch() {
	s.addStateLeaf()
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.txCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result, s.batch.Type)
	s.validateTxCommitment(result, s.txCommitment)
}

func (s *GetBatchTestSuite) TestGetBatchByID_MassMigrationBatch() {
	s.addStateLeaf()

	s.batch.Type = batchtype.MassMigration
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.mmCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result, s.batch.Type)
	s.validateMMCommitment(result, s.mmCommitment)
}

func (s *GetBatchTestSuite) TestGetBatchByID_DepositBatch() {
	s.addStateLeaf()

	s.batch.Type = batchtype.Deposit
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.depositCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result, s.batch.Type)
	s.validateDepositCommitment(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID_GenesisBatch() {
	genesisBatch := models.Batch{
		ID:                models.MakeUint256(0),
		Type:              batchtype.Genesis,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(10),
	}
	err := s.storage.AddBatch(&genesisBatch)
	s.NoError(err)

	result, err := s.api.GetBatchByID(genesisBatch.ID)
	s.NoError(err)
	s.Equal(genesisBatch.Hash, result.Hash)
	s.Equal(genesisBatch.Type, result.Type)
	s.Equal(*genesisBatch.FinalisationBlock, *result.SubmissionBlock)
	s.Nil(result.Commitments)
}

func (s *GetBatchTestSuite) TestGetBatchByID_SubmittedTxBatch() {
	s.addStateLeaf()

	batchType := batchtype.Transfer
	s.addSubmittedBatch(batchType)

	pendingCommitment := *s.txCommitment
	pendingCommitment.BodyHash = nil
	err := s.storage.AddCommitment(&pendingCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateSubmittedBatch(result, batchType)
	s.validateTxCommitment(result, &pendingCommitment)
}

func (s *GetBatchTestSuite) TestGetBatchByID_SubmittedMMBatch() {
	s.addStateLeaf()

	batchType := batchtype.MassMigration
	s.addSubmittedBatch(batchType)

	pendingCommitment := *s.mmCommitment
	pendingCommitment.BodyHash = nil
	err := s.storage.AddCommitment(&pendingCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateSubmittedBatch(result, batchType)
	s.validateMMCommitment(result, &pendingCommitment)
}

func (s *GetBatchTestSuite) TestGetBatchByID_SubmittedDepositBatch() {
	s.addStateLeaf()

	batchType := batchtype.Deposit
	s.addSubmittedBatch(batchType)

	err := s.storage.AddCommitment(s.depositCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateSubmittedBatch(result, batchType)
	s.validateDepositCommitment(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID_NonexistentBatch() {
	result, err := s.api.GetBatchByID(models.MakeUint256(0))
	s.Equal(s.batchNotFoundAPIErr, err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) addStateLeaf() {
	_, err := s.storage.StateTree.Set(uint32(1), &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *GetBatchTestSuite) addSubmittedBatch(batchType batchtype.BatchType) {
	pendingBatch := *s.batch
	pendingBatch.Type = batchType
	pendingBatch.Hash = nil
	pendingBatch.FinalisationBlock = nil
	pendingBatch.SubmissionTime = nil
	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)
}

func (s *GetBatchTestSuite) validateBatch(result *dto.BatchWithRootAndCommitments, batchType batchtype.BatchType) {
	submissionBlock := *s.batch.FinalisationBlock - config.DefaultBlocksToFinalise
	expectedBatch := dto.Batch{
		ID:                s.batch.ID,
		Hash:              s.batch.Hash,
		Type:              batchType,
		TransactionHash:   s.batch.TransactionHash,
		SubmissionBlock:   &submissionBlock,
		SubmissionTime:    s.batch.SubmissionTime,
		Status:            batchstatus.Mined,
		FinalisationBlock: s.batch.FinalisationBlock,
	}

	s.Equal(expectedBatch, result.Batch)
}

func (s *GetBatchTestSuite) validateSubmittedBatch(result *dto.BatchWithRootAndCommitments, batchType batchtype.BatchType) {
	expectedBatch := dto.Batch{
		ID:                s.batch.ID,
		Hash:              nil,
		Type:              batchType,
		TransactionHash:   s.batch.TransactionHash,
		SubmissionBlock:   nil,
		SubmissionTime:    nil,
		Status:            batchstatus.Submitted,
		FinalisationBlock: nil,
	}

	s.Equal(expectedBatch, result.Batch)
}

func (s *GetBatchTestSuite) validateTxCommitment(result *dto.BatchWithRootAndCommitments, commitment *models.TxCommitment) {
	s.Len(result.Commitments, 1)

	expectedCommitments := []dto.BatchTxCommitment{
		{
			ID: dto.CommitmentID{
				BatchID:      commitment.ID.BatchID,
				IndexInBatch: commitment.ID.IndexInBatch,
			},
			PostStateRoot:      commitment.PostStateRoot,
			LeafHash:           dto.LeafHashOrNil(commitment),
			TokenID:            models.MakeUint256(1),
			FeeReceiverStateID: commitment.FeeReceiver,
			CombinedSignature:  commitment.CombinedSignature,
		},
	}

	s.Equal(expectedCommitments, result.Commitments)
}

func (s *GetBatchTestSuite) validateMMCommitment(result *dto.BatchWithRootAndCommitments, commitment *models.MMCommitment) {
	s.Len(result.Commitments, 1)

	expectedCommitments := []dto.BatchMMCommitment{
		{
			ID: dto.CommitmentID{
				BatchID:      commitment.ID.BatchID,
				IndexInBatch: commitment.ID.IndexInBatch,
			},
			PostStateRoot:     commitment.PostStateRoot,
			LeafHash:          dto.LeafHashOrNil(commitment),
			CombinedSignature: commitment.CombinedSignature,
			WithdrawRoot:      commitment.WithdrawRoot,
			Meta: dto.MassMigrationMeta{
				SpokeID:            commitment.Meta.SpokeID,
				TokenID:            commitment.Meta.TokenID,
				Amount:             commitment.Meta.Amount,
				FeeReceiverStateID: commitment.Meta.FeeReceiver,
			},
		},
	}

	s.Equal(expectedCommitments, result.Commitments)
}

func (s *GetBatchTestSuite) validateDepositCommitment(result *dto.BatchWithRootAndCommitments) {
	s.Len(result.Commitments, 1)

	expectedCommitments := []dto.BatchDepositCommitment{
		{
			ID: dto.CommitmentID{
				BatchID:      s.depositCommitment.ID.BatchID,
				IndexInBatch: s.depositCommitment.ID.IndexInBatch,
			},
			PostStateRoot: s.depositCommitment.PostStateRoot,
			LeafHash:      dto.LeafHashOrNil(s.depositCommitment),
			SubtreeID:     s.depositCommitment.SubtreeID,
			SubtreeRoot:   s.depositCommitment.SubtreeRoot,
			Deposits:      dto.MakePendingDeposits(s.depositCommitment.Deposits),
		},
	}

	s.Equal(expectedCommitments, result.Commitments)
}

func TestGetBatchTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchTestSuite))
}
