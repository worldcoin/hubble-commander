package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
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
	transfer dto.Transfer
	wallet   *bls.Wallet
	storage  *st.TestStorage
	domain   *bls.Domain
}

func (s *GetTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransactionTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.api = &API{
		cfg:     &config.APIConfig{},
		storage: s.storage.Storage,
		client: &eth.Client{
			ChainState: chainState,
		},
		devMode: false,
	}

	err = s.storage.SetChainState(&chainState)
	s.NoError(err)
	s.domain, err = s.storage.GetDomain(chainState.ChainID)
	s.NoError(err)
	s.wallet, err = bls.NewRandomWallet(*s.domain)
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  123,
		PublicKey: *s.wallet.PublicKey(),
	})
	s.NoError(err)

	err = st.NewStateTree(s.storage.Storage).Set(1, &models.UserState{
		PubKeyID: 123,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
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
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetTransactionTestSuite) TestGetTransaction_Transfer() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	res, err := s.api.GetTransaction(*hash)
	s.NoError(err)

	transfer, ok := res.(*dto.TransferReceipt)
	s.True(ok)
	s.Equal(txstatus.Pending, transfer.Status)
}

func (s *GetTransactionTestSuite) TestGetTransaction_Create2Transfer() {
	receiverWallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  10,
		PublicKey: *receiverWallet.PublicKey(),
	})
	s.NoError(err)

	c2t := create2TransferWithoutSignature
	c2t.ToPublicKey = receiverWallet.PublicKey()
	signedTransfer, err := SignCreate2Transfer(s.wallet, c2t)
	s.NoError(err)

	hash, err := s.api.SendTransaction(dto.MakeTransaction(*signedTransfer))
	s.NoError(err)

	res, err := s.api.GetTransaction(*hash)
	s.NoError(err)

	transfer, ok := res.(*dto.Create2TransferReceipt)
	s.True(ok)
	s.Equal(txstatus.Pending, transfer.Status)
}

func TestGetTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionTestSuite))
}
