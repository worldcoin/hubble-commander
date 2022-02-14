package admin

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const authKeyValue = "secret key"

type GetPendingBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.TestStorage
	client  *eth.TestClient
	batches []models.Batch
	batch   models.Batch
}

func (s *GetPendingBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetPendingBatchesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.api = &API{
		cfg:     &config.APIConfig{AuthenticationKey: authKeyValue},
		storage: s.storage.Storage,
		client:  s.client.Client,
	}

	s.batch = models.Batch{
		ID:              models.MakeUint256(1),
		Type:            batchtype.Create2Transfer,
		TransactionHash: utils.RandomHash(),
		PrevStateRoot:   utils.NewRandomHash(),
	}

	s.batches = []models.Batch{
		{
			ID:                models.MakeUint256(1),
			Type:              batchtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			PrevStateRoot:     utils.NewRandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(42000),
			MinedTime:         models.NewTimestamp(time.Unix(140, 0).UTC()),
		},
		{
			ID:              models.MakeUint256(2),
			Type:            batchtype.Create2Transfer,
			TransactionHash: utils.RandomHash(),
			PrevStateRoot:   utils.NewRandomHash(),
		},
		{
			ID:              models.MakeUint256(3),
			Type:            batchtype.MassMigration,
			TransactionHash: utils.RandomHash(),
			PrevStateRoot:   utils.NewRandomHash(),
		},
	}
}

func (s *GetPendingBatchesTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetPendingBatchesTestSuite) TestGetPendingBatches_DifferentBatchTypes() {
	s.addBatches()

	batches, err := s.api.GetPendingBatches(contextWithAuthKey(authKeyValue))
	s.NoError(err)

	expected := s.batches[1:]
	for i := range batches {
		s.Equal(expected[i].ID, batches[i].ID)
		s.Equal(expected[i].Type, batches[i].Type)
		s.Equal(expected[i].TransactionHash, batches[i].TransactionHash)
		s.Equal(*expected[i].PrevStateRoot, batches[i].PrevStateRoot)
	}
}

func (s *GetPendingBatchesTestSuite) TestGetPendingBatches_TransferBatch() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	pendingBatch := s.addPendingTransferBatch(&tx)

	batches, err := s.api.GetPendingBatches(contextWithAuthKey(authKeyValue))
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(*pendingBatch, batches[0])
}

func (s *GetPendingBatchesTestSuite) TestGetPendingBatches_Create2TransferBatch() {
	tx := testutils.MakeCreate2Transfer(0, ref.Uint32(1), 0, 100, &models.PublicKey{1, 2, 3})
	pendingBatch := s.addPendingCT2Batch(&tx)

	batches, err := s.api.GetPendingBatches(contextWithAuthKey(authKeyValue))
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(*pendingBatch, batches[0])
}

func (s *GetPendingBatchesTestSuite) TestGetPendingBatches_MassMigrationBatch() {
	pendingBatch := s.addPendingMMBatch()

	batches, err := s.api.GetPendingBatches(contextWithAuthKey(authKeyValue))
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(*pendingBatch, batches[0])
}

func (s *GetPendingBatchesTestSuite) TestGetPendingBatches_DepositBatch() {
	pendingBatch := s.addPendingDepositBatch()

	batches, err := s.api.GetPendingBatches(contextWithAuthKey(authKeyValue))
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(*pendingBatch, batches[0])
}

func (s *GetPendingBatchesTestSuite) TestGetPendingBatches_NoBatches() {
	batches, err := s.api.GetPendingBatches(contextWithAuthKey(authKeyValue))
	s.NoError(err)
	s.Len(batches, 0)
}

func (s *GetPendingBatchesTestSuite) addBatches() {
	for i := range s.batches {
		err := s.storage.AddBatch(&s.batches[i])
		s.NoError(err)
	}
}

func (s *GetPendingBatchesTestSuite) addPendingTransferBatch(tx *models.Transfer) *dto.PendingBatch {
	batch := s.batch
	batch.Type = batchtype.BatchType(tx.Type())
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	commitment := &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type:          batch.Type,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       0,
		CombinedSignature: models.MakeRandomSignature(),
		BodyHash:          nil,
	}

	err = s.storage.AddCommitment(commitment)
	s.NoError(err)

	tx.CommitmentID = &commitment.ID
	err = s.storage.AddTransaction(tx)
	s.NoError(err)

	return newPendingBatch(&batch, commitment, models.MakeGenericArray(tx))
}

func (s *GetPendingBatchesTestSuite) addPendingCT2Batch(tx *models.Create2Transfer) *dto.PendingBatch {
	batch := s.batch
	batch.Type = batchtype.BatchType(tx.Type())
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	commitment := &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type:          batch.Type,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       0,
		CombinedSignature: models.MakeRandomSignature(),
		BodyHash:          nil,
	}

	err = s.storage.AddCommitment(commitment)
	s.NoError(err)

	tx.CommitmentID = &commitment.ID
	err = s.storage.AddTransaction(tx)
	s.NoError(err)

	return newPendingBatch(&batch, commitment, models.MakeGenericArray(tx))
}

func (s *GetPendingBatchesTestSuite) addPendingMMBatch() *dto.PendingBatch {
	batch := s.batch
	batch.Type = batchtype.MassMigration
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	commitment := &models.MMCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type:          batch.Type,
			PostStateRoot: utils.RandomHash(),
		},
		CombinedSignature: models.MakeRandomSignature(),
		BodyHash:          nil,
		Meta: &models.MassMigrationMeta{
			SpokeID:     0,
			TokenID:     models.MakeUint256(0),
			Amount:      models.MakeUint256(100),
			FeeReceiver: 0,
		},
		WithdrawRoot: utils.RandomHash(),
	}

	tx := testutils.MakeMassMigration(0, 1, 0, 100)
	s.addCommitmentWithTx(commitment, &tx)

	return newPendingBatch(&batch, commitment, models.MakeGenericArray(&tx))
}

func (s *GetPendingBatchesTestSuite) addPendingDepositBatch() *dto.PendingBatch {
	batch := s.batch
	batch.Type = batchtype.Deposit
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	commitment := &models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type:          batch.Type,
			PostStateRoot: utils.RandomHash(),
		},
		SubtreeID:   models.MakeUint256(0),
		SubtreeRoot: utils.RandomHash(),
		Deposits:    testutils.GetFourDeposits(),
	}

	err = s.storage.AddCommitment(commitment)
	s.NoError(err)

	return newPendingBatch(&batch, commitment, nil)
}

func (s *GetPendingBatchesTestSuite) addCommitmentWithTx(commitment models.Commitment, tx models.GenericTransaction) {
	err := s.storage.AddCommitment(commitment)
	s.NoError(err)

	tx.GetBase().CommitmentID = &commitment.GetCommitmentBase().ID

	err = s.storage.AddTransaction(tx)
	s.NoError(err)
}

func newPendingBatch(batch *models.Batch, commitment models.Commitment, txs models.GenericTransactionArray) *dto.PendingBatch {
	return &dto.PendingBatch{
		ID:              batch.ID,
		Type:            batch.Type,
		TransactionHash: batch.TransactionHash,
		PrevStateRoot:   *batch.PrevStateRoot,
		Commitments: []dto.PendingCommitment{
			{
				Commitment:   commitment,
				Transactions: txs,
			},
		},
	}
}

func TestGetPendingBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetPendingBatchesTestSuite))
}
