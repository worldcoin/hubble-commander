package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	api                      *API
	storage                  *st.TestStorage
	batch                    *models.Batch
	txCommitment             *models.TxCommitment
	mmCommitment             *models.MMCommitment
	depositCommitment        *models.DepositCommitment
	transfer                 models.Transfer
	create2transfer          models.Create2Transfer
	massMigration            models.MassMigration
	commitmentNotFoundAPIErr *APIError
}

func (s *GetCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage}

	s.batch = &models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(113),
		MinedTime:         models.NewTimestamp(time.Unix(140, 0).UTC()),
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
			TokenID:     models.MakeUint256(2),
			Amount:      models.MakeUint256(3),
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

	s.transfer = testutils.MakeTransfer(1, 2, 0, 100)
	s.transfer.CommitmentID = &s.txCommitment.ID

	s.create2transfer = testutils.MakeCreate2Transfer(1, ref.Uint32(2), 0, 100, &models.PublicKey{1, 2, 3})
	s.create2transfer.CommitmentID = &s.txCommitment.ID

	s.massMigration = testutils.MakeMassMigration(1, 2, 0, 100)
	s.massMigration.CommitmentID = &s.mmCommitment.ID

	s.commitmentNotFoundAPIErr = &APIError{
		Code:    20000,
		Message: "commitment not found",
	}
}

func (s *GetCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_TransferType() {
	s.addStateLeaf()
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.txCommitment)
	s.NoError(err)

	err = s.storage.AddTransaction(&s.transfer)
	s.NoError(err)

	expectedTransactions := []dto.TransferForCommitment{dto.MakeTransferForCommitment(&s.transfer)}

	commitment, err := s.api.GetCommitment(s.txCommitment.ID)
	s.NoError(err)
	s.validateTxCommitment(
		commitment.(*dto.TxCommitment),
		expectedTransactions,
		batchstatus.Mined,
		s.batch.MinedTime,
	)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_Create2TransferType() {
	s.addStateLeaf()
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	s.txCommitment.Type = batchtype.Create2Transfer
	err = s.storage.AddCommitment(s.txCommitment)
	s.NoError(err)

	err = s.storage.AddTransaction(&s.create2transfer)
	s.NoError(err)

	expectedTransactions := []dto.Create2TransferForCommitment{dto.MakeCreate2TransferForCommitment(&s.create2transfer)}

	commitment, err := s.api.GetCommitment(s.txCommitment.ID)
	s.NoError(err)
	s.validateTxCommitment(
		commitment.(*dto.TxCommitment),
		expectedTransactions,
		batchstatus.Mined,
		s.batch.MinedTime,
	)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_MassMigrationType() {
	s.addStateLeaf()
	s.batch.Type = batchtype.MassMigration
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.mmCommitment)
	s.NoError(err)

	err = s.storage.AddTransaction(&s.massMigration)
	s.NoError(err)

	expectedMassMigrations := []dto.MassMigrationForCommitment{dto.MakeMassMigrationForCommitment(&s.massMigration)}

	commitment, err := s.api.GetCommitment(s.mmCommitment.ID)
	s.NoError(err)
	s.validateMMCommitment(
		commitment.(*dto.MMCommitment),
		expectedMassMigrations,
		batchstatus.Mined,
		s.batch.MinedTime,
	)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_DepositType() {
	s.addStateLeaf()
	s.batch.Type = batchtype.Deposit
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.depositCommitment)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(s.depositCommitment.ID)
	s.NoError(err)
	s.validateDepositCommitment(
		commitment.(*dto.DepositCommitment),
		batchstatus.Mined,
		s.batch.MinedTime,
	)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_SubmittedTxCommitment() {
	s.addStateLeaf()
	s.addSubmittedBatch(batchtype.Transfer)

	pendingCommitment := *s.txCommitment
	pendingCommitment.BodyHash = nil
	err := s.storage.AddCommitment(&pendingCommitment)
	s.NoError(err)

	err = s.storage.AddTransaction(&s.transfer)
	s.NoError(err)

	expectedTransactions := []dto.TransferForCommitment{dto.MakeTransferForCommitment(&s.transfer)}

	commitment, err := s.api.GetCommitment(pendingCommitment.ID)
	s.NoError(err)
	s.validateTxCommitment(
		commitment.(*dto.TxCommitment),
		expectedTransactions,
		batchstatus.Submitted,
		nil,
	)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_SubmittedMMCommitment() {
	s.addStateLeaf()
	s.addSubmittedBatch(batchtype.MassMigration)

	pendingCommitment := *s.mmCommitment
	pendingCommitment.BodyHash = nil
	err := s.storage.AddCommitment(&pendingCommitment)
	s.NoError(err)

	err = s.storage.AddTransaction(&s.massMigration)
	s.NoError(err)

	expectedMassMigrations := []dto.MassMigrationForCommitment{dto.MakeMassMigrationForCommitment(&s.massMigration)}

	commitment, err := s.api.GetCommitment(s.mmCommitment.ID)
	s.NoError(err)
	s.validateMMCommitment(
		commitment.(*dto.MMCommitment),
		expectedMassMigrations,
		batchstatus.Submitted,
		nil,
	)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_SubmittedDepositCommitment() {
	s.addStateLeaf()
	s.addSubmittedBatch(batchtype.Deposit)

	err := s.storage.AddCommitment(s.depositCommitment)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(s.depositCommitment.ID)
	s.NoError(err)
	s.validateDepositCommitment(
		commitment.(*dto.DepositCommitment),
		batchstatus.Submitted,
		nil,
	)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_NonexistentCommitment() {
	commitment, err := s.api.GetCommitment(commitment.ID)
	s.Equal(s.commitmentNotFoundAPIErr, err)
	s.Nil(commitment)
}

func (s *GetCommitmentTestSuite) addStateLeaf() {
	_, err := s.storage.StateTree.Set(uint32(1), &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *GetCommitmentTestSuite) addSubmittedBatch(batchType batchtype.BatchType) {
	pendingBatch := *s.batch
	pendingBatch.Type = batchType
	pendingBatch.Hash = nil
	pendingBatch.FinalisationBlock = nil
	pendingBatch.MinedTime = nil
	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)
}

func (s *GetCommitmentTestSuite) validateTxCommitment(
	commitment *dto.TxCommitment,
	transactions interface{},
	batchStatus batchstatus.BatchStatus,
	batchTime *models.Timestamp,
) {
	s.NotNil(commitment)

	var leafHash *common.Hash
	if batchStatus != batchstatus.Submitted {
		leafHash = dto.LeafHashOrNil(s.txCommitment)
	}

	expectedCommitment := dto.TxCommitment{
		ID: dto.CommitmentID{
			BatchID:      s.txCommitment.ID.BatchID,
			IndexInBatch: s.txCommitment.ID.IndexInBatch,
		},
		Type:               s.txCommitment.Type,
		PostStateRoot:      s.txCommitment.PostStateRoot,
		LeafHash:           leafHash,
		TokenID:            models.MakeUint256(1),
		FeeReceiverStateID: s.txCommitment.FeeReceiver,
		CombinedSignature:  s.txCommitment.CombinedSignature,
		Status:             batchStatus,
		BatchTime:          batchTime,
		Transactions:       transactions,
	}

	s.Equal(expectedCommitment, *commitment)
}

func (s *GetCommitmentTestSuite) validateMMCommitment(
	commitment *dto.MMCommitment,
	transactions interface{},
	batchStatus batchstatus.BatchStatus,
	batchTime *models.Timestamp,
) {
	s.NotNil(commitment)

	var leafHash *common.Hash
	if batchStatus != batchstatus.Submitted {
		leafHash = dto.LeafHashOrNil(s.mmCommitment)
	}

	expectedCommitment := dto.MMCommitment{
		ID: dto.CommitmentID{
			BatchID:      s.mmCommitment.ID.BatchID,
			IndexInBatch: s.mmCommitment.ID.IndexInBatch,
		},
		Type:              s.mmCommitment.Type,
		PostStateRoot:     s.mmCommitment.PostStateRoot,
		LeafHash:          leafHash,
		CombinedSignature: s.mmCommitment.CombinedSignature,
		Status:            batchStatus,
		BatchTime:         batchTime,
		WithdrawRoot:      s.mmCommitment.WithdrawRoot,
		Meta: dto.MassMigrationMeta{
			SpokeID:            s.mmCommitment.Meta.SpokeID,
			TokenID:            s.mmCommitment.Meta.TokenID,
			Amount:             s.mmCommitment.Meta.Amount,
			FeeReceiverStateID: s.mmCommitment.Meta.FeeReceiver,
		},
		Transactions: transactions,
	}

	s.Equal(expectedCommitment, *commitment)
}

func (s *GetCommitmentTestSuite) validateDepositCommitment(
	commitment *dto.DepositCommitment,
	batchStatus batchstatus.BatchStatus,
	batchTime *models.Timestamp,
) {
	s.NotNil(commitment)

	leafHash := s.depositCommitment.LeafHash()
	expectedCommitment := dto.DepositCommitment{
		ID: dto.CommitmentID{
			BatchID:      s.depositCommitment.ID.BatchID,
			IndexInBatch: s.depositCommitment.ID.IndexInBatch,
		},
		Type:          s.depositCommitment.Type,
		PostStateRoot: s.depositCommitment.PostStateRoot,
		LeafHash:      &leafHash,
		Status:        batchStatus,
		BatchTime:     batchTime,
		SubtreeID:     s.depositCommitment.SubtreeID,
		SubtreeRoot:   s.depositCommitment.SubtreeRoot,
		Deposits:      dto.MakeDeposits(s.depositCommitment.Deposits),
	}

	s.Equal(expectedCommitment, *commitment)
}

func TestGetCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(GetCommitmentTestSuite))
}
