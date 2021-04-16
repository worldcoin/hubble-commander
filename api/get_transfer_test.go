package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransferTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	db       *db.TestDB
	transfer dto.Transfer
	wallet   *bls.Wallet
}

func (s *GetTransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransferTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)
	s.api = &API{nil, storage, nil}
	s.db = testDB

	s.wallet, err = bls.NewRandomWallet(mockDomain)
	s.NoError(err)

	err = storage.AddAccountIfNotExists(&models.Account{
		AccountIndex: 123,
		PublicKey:    *s.wallet.PublicKey(),
	})
	s.NoError(err)

	err = st.NewStateTree(storage).Set(1, &userState)
	s.NoError(err)

	s.transfer = s.signTransfer(transferWithoutSignature)
}

func (s *GetTransferTestSuite) signTransfer(transfer dto.Transfer) dto.Transfer {
	signedTransfer, err := SignTransfer(s.wallet, transfer)
	s.NoError(err)
	return *signedTransfer
}

func (s *GetTransferTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetTransferTestSuite) TestApi_GetTransfer() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	res, err := s.api.GetTransfer(*hash)
	s.NoError(err)

	s.Equal(models.Pending, res.Status)
}

func TestGetTransferTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransferTestSuite))
}
