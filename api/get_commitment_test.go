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
	batch                    models.Batch
	commitment               models.TxCommitment
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

	s.batch = models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(113),
		SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
	}

	s.commitment = commitment
	s.commitment.ID.BatchID = s.batch.ID
	s.commitment.ID.IndexInBatch = 0

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
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddTxCommitment(&s.commitment)
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
			CommitmentID: &s.commitment.ID,
		},
		ToStateID: 2,
	}
	err = s.storage.AddTransfer(&transfer)
	s.NoError(err)

	expectedCommitment := &dto.Commitment{
		ID:                *dto.MakeCommitmentID(&s.commitment.ID),
		Type:              s.commitment.Type,
		PostStateRoot:     s.commitment.PostStateRoot,
		FeeReceiver:       s.commitment.FeeReceiver,
		CombinedSignature: s.commitment.CombinedSignature,
		BodyHash:          s.commitment.BodyHash,
		Status:            txstatus.InBatch,
		BatchTime:         s.batch.SubmissionTime,
		Transactions: []dto.TransferForCommitment{{
			Hash:        transfer.Hash,
			FromStateID: transfer.FromStateID,
			Amount:      transfer.Amount,
			Fee:         transfer.Fee,
			Nonce:       transfer.Nonce,
			Signature:   transfer.Signature,
			ReceiveTime: transfer.ReceiveTime,
			ToStateID:   transfer.ToStateID,
		}},
	}

	commitment, err := s.api.GetCommitment(s.commitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitment, commitment)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_Create2TransferType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	s.commitment.Type = batchtype.Create2Transfer
	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Create2Transfer,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			CommitmentID: &s.commitment.ID,
		},
		ToStateID:   ref.Uint32(2),
		ToPublicKey: models.PublicKey{2, 3, 4},
	}
	err = s.storage.AddCreate2Transfer(&transfer)
	s.NoError(err)

	expectedCommitment := &dto.Commitment{
		ID:                *dto.MakeCommitmentID(&s.commitment.ID),
		Type:              s.commitment.Type,
		PostStateRoot:     s.commitment.PostStateRoot,
		FeeReceiver:       s.commitment.FeeReceiver,
		CombinedSignature: s.commitment.CombinedSignature,
		BodyHash:          s.commitment.BodyHash,
		Status:            txstatus.InBatch,
		BatchTime:         s.batch.SubmissionTime,
		Transactions: []dto.Create2TransferForCommitment{{
			Hash:        transfer.Hash,
			FromStateID: transfer.FromStateID,
			Amount:      transfer.Amount,
			Fee:         transfer.Fee,
			Nonce:       transfer.Nonce,
			Signature:   transfer.Signature,
			ReceiveTime: transfer.ReceiveTime,
			ToStateID:   transfer.ToStateID,
			ToPublicKey: transfer.ToPublicKey,
		}},
	}

	commitment, err := s.api.GetCommitment(s.commitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitment, commitment)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_MassMigrationType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	s.commitment.Type = batchtype.MassMigration
	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	massMigration := models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.MassMigration,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			CommitmentID: &s.commitment.ID,
		},
		SpokeID: 2,
	}
	err = s.storage.AddMassMigration(&massMigration)
	s.NoError(err)

	expectedCommitment := &dto.Commitment{
		ID:                *dto.MakeCommitmentID(&s.commitment.ID),
		Type:              s.commitment.Type,
		PostStateRoot:     s.commitment.PostStateRoot,
		FeeReceiver:       s.commitment.FeeReceiver,
		CombinedSignature: s.commitment.CombinedSignature,
		BodyHash:          s.commitment.BodyHash,
		Status:            txstatus.InBatch,
		BatchTime:         s.batch.SubmissionTime,
		Transactions: []dto.MassMigrationForCommitment{{
			Hash:        massMigration.Hash,
			FromStateID: massMigration.FromStateID,
			Amount:      massMigration.Amount,
			Fee:         massMigration.Fee,
			Nonce:       massMigration.Nonce,
			Signature:   massMigration.Signature,
			ReceiveTime: massMigration.ReceiveTime,
			SpokeID:     massMigration.SpokeID,
		}},
	}

	commitment, err := s.api.GetCommitment(s.commitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitment, commitment)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_PendingBatch() {
	pendingBatch := s.batch
	pendingBatch.Hash = nil
	pendingBatch.FinalisationBlock = nil
	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Transfer,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			CommitmentID: &s.commitment.ID,
		},
		ToStateID: 2,
	}
	err = s.storage.AddTransfer(&transfer)
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

func TestGetCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(GetCommitmentTestSuite))
}
