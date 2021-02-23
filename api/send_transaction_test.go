package api

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SendTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	api *API
	db  *db.TestDB
}

func (s *SendTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SendTransactionTestSuite) SetupTest() {
	testDB, err := db.GetTestDB()
	s.NoError(err)

	storage := &st.Storage{DB: testDB.DB}
	s.api = &API{nil, storage}
	s.db = testDB
}

func (s *SendTransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction() {
	tx := models.IncomingTransaction{
		FromIndex: big.NewInt(1),
		ToIndex:   big.NewInt(2),
		Amount:    big.NewInt(50),
		Fee:       big.NewInt(10),
		Nonce:     big.NewInt(0),
		Signature: []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}
	hash, err := s.api.SendTransaction(tx)
	s.NoError(err)
	s.Equal(common.HexToHash("0x3e136a19201d6fc73c4e3c76951edfb94eb9a7a0c7e15492696ffddb3e1b2c68"), *hash)
}

func TestSendTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransactionTestSuite))
}
