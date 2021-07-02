package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
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
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage}

	s.batch = models.Batch{
		ID:                models.MakeUint256(1),
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(113),
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
	s.NoError(err)
	s.NotNil(commitment)
	s.Len(commitment.Transactions, 1)
}

func (s *GetCommitmentTestSuite) TestGetCommitment_Create2TransferType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	c2tCommitment := s.commitment
	c2tCommitment.Type = txtype.Create2Transfer
	commitmentID, err := s.storage.AddCommitment(&c2tCommitment)
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  2,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	create2Transfer := models.Create2Transfer{
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
	_, err = s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(*commitmentID)
	s.NoError(err)
	s.NotNil(commitment)
	s.Len(commitment.Transactions, 1)
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
