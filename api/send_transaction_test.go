package api

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SendTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	cfg     *config.Config
	storage db.Storage
}

func (s *SendTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SendTransactionTestSuite) SetupTest() {
	cfg := config.GetTestConfig()
	s.cfg = &cfg
	storage, err := db.GetTestStorage()
	s.NoError(err)
	s.storage = storage
}

func (s *SendTransactionTestSuite) TearDownTest() {
	migrator, err := db.GetMigrator(s.cfg)
	s.NoError(err)

	s.NoError(migrator.Down())

	err = s.storage.DB.Close()
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction() {
	api := Api{s.cfg, &s.storage}
	tx := models.IncomingTransaction{
		FromIndex: big.NewInt(1),
		ToIndex:   big.NewInt(2),
		Amount:    big.NewInt(50),
		Fee:       big.NewInt(10),
		Nonce:     big.NewInt(0),
		Signature: []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}
	hash, err := api.SendTransaction(tx)
	s.NoError(err)
	s.Equal(common.HexToHash("0x3e136a19201d6fc73c4e3c76951edfb94eb9a7a0c7e15492696ffddb3e1b2c68"), *hash)
}

func TestSendTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransactionTestSuite))
}
