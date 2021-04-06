package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	api *API
	db  *db.TestDB
	tx  *models.Transaction
}

func (s *GetTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)
	s.api = &API{nil, storage, nil}
	s.db = testDB

	userState := models.UserState{
		AccountIndex: 1,
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}

	tree := st.NewStateTree(storage)
	err = tree.Set(1, &userState)
	s.NoError(err)

	tx := &models.Transaction{
		FromIndex: 1,
		ToIndex:   2,
		Amount:    *models.NewUint256(50),
		Fee:       *models.NewUint256(10),
		Nonce:     *models.NewUint256(0),
		Signature: []byte{1, 2, 3, 4},
	}

	s.tx = tx
}

func (s *GetTransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetTransactionTestSuite) TestApi_GetTransaction() {
	incomingTx := models.IncomingTransaction{
		FromIndex: models.NewUint256(1),
		ToIndex:   models.NewUint256(2),
		Amount:    models.NewUint256(50),
		Fee:       models.NewUint256(10),
		Nonce:     models.NewUint256(0),
		Signature: []byte{1, 2, 3, 4},
	}

	hash, err := s.api.SendTransaction(incomingTx)
	s.NoError(err)

	res, err := s.api.GetTransaction(*hash)
	s.NoError(err)

	s.Equal(models.Pending, res.Status)
}

func TestGetTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionTestSuite))
}
