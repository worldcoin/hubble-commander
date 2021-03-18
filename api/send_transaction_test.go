package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SendTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	api  *API
	db   *db.TestDB
	tree *storage.StateTree
}

func (s *SendTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SendTransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)
	s.tree = st.NewStateTree(storage)
	s.api = &API{nil, storage}
	s.db = testDB

	userState := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}

	err = s.tree.Set(1, &userState)
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction() {
	tx := models.IncomingTransaction{
		FromIndex: models.NewUint256(1),
		ToIndex:   models.NewUint256(2),
		Amount:    models.NewUint256(50),
		Fee:       models.NewUint256(10),
		Nonce:     models.NewUint256(0),
		Signature: []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}
	hash, err := s.api.SendTransaction(tx)
	s.NoError(err)
	s.Equal(common.HexToHash("0x0757e6b9b057336b010007d26489dc3a46d89a5349824965b9129ca26ff72340"), *hash)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_VerifyNonce_TooHigh() {
	tx := models.IncomingTransaction{
		FromIndex: models.NewUint256(1),
		ToIndex:   models.NewUint256(2),
		Amount:    models.NewUint256(50),
		Fee:       models.NewUint256(10),
		Nonce:     models.NewUint256(10),
		Signature: []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}
	_, err := s.api.SendTransaction(tx)
	s.EqualErrorf(err, "nonce too high", err.Error())
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_VerifyNonce_TooLow() {
	userState := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(1),
	}

	err := s.tree.Set(2, &userState)
	s.NoError(err)

	tx := models.IncomingTransaction{
		FromIndex: models.NewUint256(2),
		ToIndex:   models.NewUint256(1),
		Amount:    models.NewUint256(50),
		Fee:       models.NewUint256(10),
		Nonce:     models.NewUint256(0),
		Signature: []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}
	_, err = s.api.SendTransaction(tx)
	s.EqualErrorf(err, "nonce too low", err.Error())
}

func TestSendTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransactionTestSuite))
}
