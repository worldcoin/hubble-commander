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

var (
	userState = models.UserState{
		AccountIndex: 1,
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	transfer = dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   utils.RandomBytes(12),
	}
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

	err = s.tree.Set(1, &userState)
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ReturnsNonNilHash() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidateNonce_TooLow() {
	userStateWithIncreasedNonce := userState
	userStateWithIncreasedNonce.Nonce = *models.NewUint256(1)

	err := s.tree.Set(1, &userStateWithIncreasedNonce)
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(ErrNonceTooLow, err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidateFee() {
	transferWithZeroFee := transfer
	transferWithZeroFee.Fee = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroFee))
	s.Equal(ErrFeeTooLow, err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidateBalance() {
	transferWithHugeAmount := transfer
	transferWithHugeAmount.Amount = models.NewUint256(500)
	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithHugeAmount))
	s.Equal(ErrNotEnoughBalance, err)
}

func TestSendTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransactionTestSuite))
}
