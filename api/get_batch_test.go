package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
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
	txCommitment        models.TxCommitment
	mmCommitment        models.MMCommitment
	depositCommitment   models.DepositCommitment
	batch               models.Batch
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

	s.batch = models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(42000),
		AccountTreeRoot:   utils.NewRandomHash(),
		SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
	}

	s.txCommitment = commitment
	s.txCommitment.ID.BatchID = s.batch.ID
	s.txCommitment.BodyHash = utils.NewRandomHash()

	s.mmCommitment = models.MMCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      s.batch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.MassMigration,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       1,
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

	s.depositCommitment = models.DepositCommitment{
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
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.txCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result)
	s.validateTxCommitment(result)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_MassMigrationBatch() {
	s.addStateLeaf()

	s.batch.Type = batchtype.MassMigration
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.mmCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result)
	s.validateMMCommitment(result)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_DepositBatch() {
	s.addStateLeaf()

	s.batch.Type = batchtype.Deposit
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.depositCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result)
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
	s.Equal(*genesisBatch.FinalisationBlock, result.SubmissionBlock)
	s.Len(result.Commitments, 0)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NonexistentBatch() {
	result, err := s.api.GetBatchByHash(utils.RandomHash())
	s.Equal(s.batchNotFoundAPIErr, err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID_TxBatch() {
	s.addStateLeaf()
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.txCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result)
	s.validateTxCommitment(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID_MassMigrationBatch() {
	s.addStateLeaf()

	s.batch.Type = batchtype.MassMigration
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.mmCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result)
	s.validateMMCommitment(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID_DepositBatch() {
	s.addStateLeaf()

	s.batch.Type = batchtype.Deposit
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.depositCommitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.validateBatch(result)
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
	s.Equal(*genesisBatch.FinalisationBlock, result.SubmissionBlock)
	s.Len(result.Commitments, 0)
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

func (s *GetBatchTestSuite) validateBatch(result *dto.BatchWithRootAndCommitments) {
	s.Equal(s.batch.ID, result.ID)
	s.Equal(s.batch.Hash, result.Hash)
	s.Equal(s.batch.Type, result.Type)
	s.Equal(s.batch.TransactionHash, result.TransactionHash)
	s.Equal(*s.batch.FinalisationBlock-config.DefaultBlocksToFinalise, result.SubmissionBlock)
	s.Equal(s.batch.FinalisationBlock, result.FinalisationBlock)
	s.Equal(s.batch.SubmissionTime, result.SubmissionTime)
}

func (s *GetBatchTestSuite) validateTxCommitment(result *dto.BatchWithRootAndCommitments) {
	s.Len(result.Commitments, 1)

	commitment := result.Commitments[0]
	s.Equal(*dto.NewCommitmentID(&s.txCommitment.ID), commitment.ID)
	s.Equal(s.txCommitment.PostStateRoot, commitment.PostStateRoot)
	s.Equal(s.txCommitment.LeafHash(), commitment.LeafHash)
	s.Equal(s.txCommitment.FeeReceiver, *commitment.FeeReceiverStateID)
	s.Equal(s.txCommitment.CombinedSignature, *commitment.CombinedSignature)

	stateLeaf, err := s.storage.StateTree.Leaf(s.txCommitment.FeeReceiver)
	s.NoError(err)

	s.Equal(stateLeaf.TokenID, *commitment.TokenID)

	s.Nil(commitment.Meta)
	s.Nil(commitment.WithdrawRoot)
	s.Nil(commitment.SubtreeID)
	s.Nil(commitment.SubtreeRoot)
	s.Nil(commitment.Deposits)
}

func (s *GetBatchTestSuite) validateMMCommitment(result *dto.BatchWithRootAndCommitments) {
	s.Len(result.Commitments, 1)

	commitment := result.Commitments[0]
	s.Equal(*dto.NewCommitmentID(&s.mmCommitment.ID), commitment.ID)
	s.Equal(s.mmCommitment.PostStateRoot, commitment.PostStateRoot)
	s.Equal(s.mmCommitment.LeafHash(), commitment.LeafHash)
	s.Nil(commitment.TokenID)
	s.Nil(commitment.FeeReceiverStateID)
	s.Equal(s.mmCommitment.CombinedSignature, *commitment.CombinedSignature)

	expectedMeta := &dto.MassMigrationMeta{
		SpokeID:            s.mmCommitment.Meta.SpokeID,
		TokenID:            s.mmCommitment.Meta.TokenID,
		Amount:             s.mmCommitment.Meta.Amount,
		FeeReceiverStateID: s.mmCommitment.Meta.FeeReceiver,
	}

	s.Equal(expectedMeta, commitment.Meta)
	s.Equal(s.mmCommitment.WithdrawRoot, *commitment.WithdrawRoot)

	s.Nil(commitment.SubtreeID)
	s.Nil(commitment.SubtreeRoot)
	s.Nil(commitment.Deposits)
}

func (s *GetBatchTestSuite) validateDepositCommitment(result *dto.BatchWithRootAndCommitments) {
	s.Len(result.Commitments, 1)

	commitment := result.Commitments[0]
	s.Equal(*dto.NewCommitmentID(&s.depositCommitment.ID), commitment.ID)
	s.Equal(s.depositCommitment.PostStateRoot, commitment.PostStateRoot)
	s.Nil(commitment.TokenID)
	s.Nil(commitment.FeeReceiverStateID)
	s.Nil(commitment.CombinedSignature)
	s.Equal(s.depositCommitment.SubtreeID, *commitment.SubtreeID)
	s.Equal(s.depositCommitment.SubtreeRoot, *commitment.SubtreeRoot)
	s.Len(s.depositCommitment.Deposits, 1)

	expectedDeposit := dto.PendingDeposit{
		ID: dto.DepositID{
			SubtreeID:    s.depositCommitment.Deposits[0].ID.SubtreeID,
			DepositIndex: s.depositCommitment.Deposits[0].ID.DepositIndex,
		},
		ToPubKeyID: s.depositCommitment.Deposits[0].ToPubKeyID,
		TokenID:    s.depositCommitment.Deposits[0].TokenID,
		L2Amount:   s.depositCommitment.Deposits[0].L2Amount,
	}
	s.Equal(expectedDeposit, commitment.Deposits[0])

	s.Nil(commitment.Meta)
	s.Nil(commitment.WithdrawRoot)
}

func TestGetBatchTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchTestSuite))
}
