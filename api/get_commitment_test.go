package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
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
	api        *API
	storage    *st.TestStorage
	batch      models.Batch
	commitment models.Commitment
}

func (s *GetCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage}

	s.batch = models.Batch{
		ID:                models.MakeUint256(1),
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(113),
		SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
	}

	s.commitment = commitment
	s.commitment.IncludedInBatch = &s.batch.ID
}

func (s *GetCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_TransferType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	commitmentID, err := s.storage.AddCommitment(&s.commitment)
	s.NoError(err)
	s.commitment.ID = *commitmentID

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 utils.RandomHash(),
			TxType:               txtype.Transfer,
			FromStateID:          1,
			Amount:               models.MakeUint256(50),
			Fee:                  models.MakeUint256(10),
			Nonce:                models.MakeUint256(0),
			Signature:            models.MakeRandomSignature(),
			IncludedInCommitment: commitmentID,
		},
		ToStateID: 2,
	}
	receiveTime, err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	expectedCommitment := &dto.Commitment{
		Commitment: s.commitment,
		Status:     txstatus.InBatch,
		BatchTime:  s.batch.SubmissionTime,
		Transactions: []models.TransferForCommitment{{
			TransactionBaseForCommitment: models.TransactionBaseForCommitment{
				Hash:        transfer.Hash,
				FromStateID: transfer.FromStateID,
				Amount:      transfer.Amount,
				Fee:         transfer.Fee,
				Nonce:       transfer.Nonce,
				Signature:   transfer.Signature,
				ReceiveTime: receiveTime,
			},
			ToStateID: transfer.ToStateID,
		}},
	}

	commitment, err := s.api.GetCommitment(*commitmentID)
	s.NoError(err)
	s.Equal(expectedCommitment, commitment)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_Create2TransferType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	s.commitment.Type = txtype.Create2Transfer
	commitmentID, err := s.storage.AddCommitment(&s.commitment)
	s.NoError(err)
	s.commitment.ID = *commitmentID

	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 utils.RandomHash(),
			TxType:               txtype.Create2Transfer,
			FromStateID:          1,
			Amount:               models.MakeUint256(50),
			Fee:                  models.MakeUint256(10),
			Nonce:                models.MakeUint256(0),
			IncludedInCommitment: commitmentID,
		},
		ToStateID:   ref.Uint32(2),
		ToPublicKey: models.PublicKey{2, 3, 4},
	}
	receiveTime, err := s.storage.AddCreate2Transfer(&transfer)
	s.NoError(err)

	expectedCommitment := &dto.Commitment{
		Commitment: s.commitment,
		Status:     txstatus.InBatch,
		BatchTime:  s.batch.SubmissionTime,
		Transactions: []models.Create2TransferForCommitment{{
			TransactionBaseForCommitment: models.TransactionBaseForCommitment{
				Hash:        transfer.Hash,
				FromStateID: transfer.FromStateID,
				Amount:      transfer.Amount,
				Fee:         transfer.Fee,
				Nonce:       transfer.Nonce,
				Signature:   transfer.Signature,
				ReceiveTime: receiveTime,
			},
			ToStateID:   transfer.ToStateID,
			ToPublicKey: transfer.ToPublicKey,
		}},
	}

	commitment, err := s.api.GetCommitment(*commitmentID)
	s.NoError(err)
	s.Equal(expectedCommitment, commitment)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_PendingBatch() {
	pendingBatch := s.batch
	pendingBatch.Hash = nil
	pendingBatch.FinalisationBlock = nil
	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	commitmentID, err := s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 utils.RandomHash(),
			TxType:               txtype.Transfer,
			FromStateID:          1,
			Amount:               models.MakeUint256(50),
			Fee:                  models.MakeUint256(10),
			Nonce:                models.MakeUint256(0),
			IncludedInCommitment: commitmentID,
		},
		ToStateID: 2,
	}
	_, err = s.storage.AddTransfer(&transfer)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(*commitmentID)
	s.Equal(st.NewNotFoundError("commitment"), err)
	s.Nil(commitment)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_NotExistingCommitment() {
	commitment, err := s.api.GetCommitment(123)
	s.Equal(st.NewNotFoundError("commitment"), err)
	s.Nil(commitment)
}

func TestGetCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(GetCommitmentTestSuite))
}
