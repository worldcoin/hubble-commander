package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SendTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	api  *API
	db   *db.TestDB
	tree *st.StateTree
}

func (s *SendTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SendTransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)
	s.tree = st.NewStateTree(storage)
	s.api = &API{nil, storage, nil}
	s.db = testDB

	userState := models.UserState{
		PubkeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}

	err = s.tree.Set(1, &userState)
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ReturnsNonNilHash() {
	transfer := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   utils.RandomBytes(12),
	}
	hash, err := s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidateNonce_TooLow() {
	userState := models.UserState{
		PubkeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(1),
	}

	err := s.tree.Set(2, &userState)
	s.NoError(err)

	transfer := dto.Transfer{
		FromStateID: ref.Uint32(2),
		ToStateID:   ref.Uint32(1),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   utils.RandomBytes(12),
	}
	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(ErrNonceTooLow, err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidateFee() {
	transfer := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(0),
		Nonce:       models.NewUint256(0),
		Signature:   utils.RandomBytes(12),
	}
	_, err := s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(ErrFeeTooLow, err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidateBalance() {
	transfer := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(500),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   utils.RandomBytes(12),
	}
	_, err := s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(ErrNotEnoughBalance, err)
}

func TestSendTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransactionTestSuite))
}
