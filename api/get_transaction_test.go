package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	db       *db.TestDB
	transfer dto.Transfer
	wallet   *bls.Wallet
}

func (s *GetTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)
	s.api = &API{
		cfg:     &config.APIConfig{},
		storage: storage,
		client:  nil,
	}
	s.db = testDB

	s.wallet, err = bls.NewRandomWallet(mockDomain)
	s.NoError(err)

	err = storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  123,
		PublicKey: *s.wallet.PublicKey(),
	})
	s.NoError(err)

	err = st.NewStateTree(storage).Set(1, &models.UserState{
		PubKeyID:   123,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)

	s.transfer = s.signTransfer(transferWithoutSignature)
}

func (s *GetTransactionTestSuite) signTransfer(transfer dto.Transfer) dto.Transfer {
	signedTransfer, err := SignTransfer(s.wallet, transfer)
	s.NoError(err)
	return *signedTransfer
}

func (s *GetTransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetTransactionTestSuite) TestGetTransaction() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	res, err := s.api.GetTransaction(*hash)
	s.NoError(err)

	transfer, ok := res.(*dto.TransferReceipt)
	s.True(ok)
	s.Equal(txstatus.Pending, transfer.Status)
}

func TestGetTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionTestSuite))
}
