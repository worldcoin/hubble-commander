package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
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
		FeeReceiver:       1,
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

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Transfer,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			Signature:    models.MakeRandomSignature(),
			CommitmentID: &s.txCommitment.ID,
		},
		ToStateID: 2,
	}
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(s.txCommitment.ID)
	s.NoError(err)
	s.NotNil(commitment)
	s.validateTxCommitment(commitment)

	expectedTransactions := []dto.TransferForCommitment{{
		Hash:        transfer.Hash,
		FromStateID: transfer.FromStateID,
		Amount:      transfer.Amount,
		Fee:         transfer.Fee,
		Nonce:       transfer.Nonce,
		Signature:   transfer.Signature,
		ReceiveTime: transfer.ReceiveTime,
		ToStateID:   transfer.ToStateID,
	}}

	s.Equal(expectedTransactions, commitment.Transactions)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_Create2TransferType() {
	s.addStateLeaf()
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	s.txCommitment.Type = batchtype.Create2Transfer
	err = s.storage.AddCommitment(s.txCommitment)
	s.NoError(err)

	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Create2Transfer,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			CommitmentID: &s.txCommitment.ID,
		},
		ToStateID:   ref.Uint32(2),
		ToPublicKey: models.PublicKey{2, 3, 4},
	}
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(s.txCommitment.ID)
	s.NoError(err)
	s.NotNil(commitment)
	s.validateTxCommitment(commitment)

	expectedTransactions := []dto.Create2TransferForCommitment{{
		Hash:        transfer.Hash,
		FromStateID: transfer.FromStateID,
		Amount:      transfer.Amount,
		Fee:         transfer.Fee,
		Nonce:       transfer.Nonce,
		Signature:   transfer.Signature,
		ReceiveTime: transfer.ReceiveTime,
		ToStateID:   transfer.ToStateID,
		ToPublicKey: transfer.ToPublicKey,
	}}

	s.Equal(expectedTransactions, commitment.Transactions)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_MassMigrationType() {
	s.addStateLeaf()
	s.batch.Type = batchtype.MassMigration
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.mmCommitment)
	s.NoError(err)

	massMigration := models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.MassMigration,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			CommitmentID: &s.mmCommitment.ID,
		},
		SpokeID: 2,
	}
	err = s.storage.AddTransaction(&massMigration)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(s.mmCommitment.ID)
	s.NoError(err)
	s.NotNil(commitment)
	s.validateMMCommitment(commitment)

	expectedMassMigrations := []dto.MassMigrationForCommitment{{
		Hash:        massMigration.Hash,
		FromStateID: massMigration.FromStateID,
		Amount:      massMigration.Amount,
		Fee:         massMigration.Fee,
		Nonce:       massMigration.Nonce,
		Signature:   massMigration.Signature,
		ReceiveTime: massMigration.ReceiveTime,
		SpokeID:     massMigration.SpokeID,
	}}

	s.Equal(expectedMassMigrations, commitment.Transactions)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_PendingBatch() {
	s.addStateLeaf()
	pendingBatch := *s.batch
	pendingBatch.Hash = nil
	pendingBatch.FinalisationBlock = nil
	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	err = s.storage.AddCommitment(s.txCommitment)
	s.NoError(err)

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Transfer,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			CommitmentID: &s.txCommitment.ID,
		},
		ToStateID: 2,
	}
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(commitment.ID)
	s.Equal(s.commitmentNotFoundAPIErr, err)
	s.Nil(commitment)
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

func (s *GetCommitmentTestSuite) validateTxCommitment(commitment *dto.Commitment) {
	s.Equal(*dto.NewCommitmentID(&s.txCommitment.ID), commitment.ID)
	s.Equal(s.txCommitment.Type, commitment.Type)
	s.Equal(s.txCommitment.PostStateRoot, commitment.PostStateRoot)
	s.Equal(s.txCommitment.LeafHash(), commitment.LeafHash)
	s.Equal(s.txCommitment.FeeReceiver, commitment.FeeReceiverStateID)
	s.Equal(s.txCommitment.CombinedSignature, commitment.CombinedSignature)
	s.Equal(txstatus.InBatch, commitment.Status)
	s.Equal(s.batch.SubmissionTime, commitment.BatchTime)

	stateLeaf, err := s.storage.StateTree.Leaf(s.txCommitment.FeeReceiver)
	s.NoError(err)

	s.Equal(stateLeaf.TokenID, commitment.TokenID)
}

func (s *GetCommitmentTestSuite) validateMMCommitment(commitment *dto.Commitment) {
	s.Equal(*dto.NewCommitmentID(&s.mmCommitment.ID), commitment.ID)
	s.Equal(s.mmCommitment.Type, commitment.Type)
	s.Equal(s.mmCommitment.PostStateRoot, commitment.PostStateRoot)
	s.Equal(s.mmCommitment.LeafHash(), commitment.LeafHash)
	s.Equal(s.mmCommitment.FeeReceiver, commitment.FeeReceiverStateID)
	s.Equal(s.mmCommitment.CombinedSignature, commitment.CombinedSignature)
	s.Equal(txstatus.InBatch, commitment.Status)
	s.Equal(s.batch.SubmissionTime, commitment.BatchTime)

	stateLeaf, err := s.storage.StateTree.Leaf(s.mmCommitment.FeeReceiver)
	s.NoError(err)

	s.Equal(stateLeaf.TokenID, commitment.TokenID)

	expectedMeta := &dto.MassMigrationMeta{
		SpokeID:     s.mmCommitment.Meta.SpokeID,
		TokenID:     s.mmCommitment.Meta.TokenID,
		Amount:      s.mmCommitment.Meta.Amount,
		FeeReceiver: s.mmCommitment.Meta.FeeReceiver,
	}

	s.Equal(expectedMeta, commitment.Meta)
}

func TestGetCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(GetCommitmentTestSuite))
}
