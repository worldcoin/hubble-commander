package api

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	api             *API
	transfer        dto.Transfer
	create2Transfer dto.Create2Transfer
	massMigration   dto.MassMigration
	wallet          *bls.Wallet
	storage         *st.TestStorage
	domain          *bls.Domain
}

func (s *GetTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransactionTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.api = &API{
		cfg:                     &config.APIConfig{},
		storage:                 s.storage.Storage,
		client:                  eth.DomainOnlyTestClient,
		commanderMetrics:        metrics.NewCommanderMetrics(),
		txPool:                  mempool.NewTestTxPool(),
		disableSignatures:       false,
		isAcceptingTransactions: true,
	}

	s.domain, err = s.api.client.GetDomain()
	s.NoError(err)
	s.wallet, err = bls.NewRandomWallet(*s.domain)
	s.NoError(err)

	err = s.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  123,
		PublicKey: *s.wallet.PublicKey(),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 123,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	s.transfer = s.signTransfer(transferWithoutSignature)
	s.create2Transfer = s.signCreate2Transfer(create2TransferWithoutSignature)
	s.massMigration = s.signMassMigration(massMigrationWithoutSignature)

	err = s.storage.AddRegisteredSpoke(&models.RegisteredSpoke{
		ID: models.MakeUint256(2),
	})
	s.NoError(err)
}

func (s *GetTransactionTestSuite) signTransfer(transfer dto.Transfer) dto.Transfer {
	signedTransfer, err := SignTransfer(s.wallet, transfer)
	s.NoError(err)
	return *signedTransfer
}

func (s *GetTransactionTestSuite) signCreate2Transfer(create2Transfer dto.Create2Transfer) dto.Create2Transfer {
	randomWallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	create2Transfer.ToPublicKey = randomWallet.PublicKey()
	signedCreate2Transfer, err := SignCreate2Transfer(s.wallet, create2Transfer)
	s.NoError(err)
	return *signedCreate2Transfer
}

func (s *GetTransactionTestSuite) signMassMigration(massMigration dto.MassMigration) dto.MassMigration {
	signedMassMigration, err := SignMassMigration(s.wallet, massMigration)
	s.NoError(err)
	return *signedMassMigration
}

func (s *GetTransactionTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetTransactionTestSuite) TestGetTransaction_Transfer() {
	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(s.transfer))
	s.NoError(err)

	receipt, err := s.api.GetTransaction(*hash)
	s.NoError(err)
	s.Equal(txstatus.Pending, receipt.Status)

	s.Equal(txtype.Transfer, receipt.TxType)
}

func (s *GetTransactionTestSuite) TestGetTransaction_Create2Transfer() {
	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(s.create2Transfer))
	s.NoError(err)

	receipt, err := s.api.GetTransaction(*hash)
	s.NoError(err)
	s.Equal(txstatus.Pending, receipt.Status)

	s.Equal(txtype.Create2Transfer, receipt.TxType)
}

func (s *GetTransactionTestSuite) TestGetTransaction_MassMigration() {
	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(s.massMigration))
	s.NoError(err)

	receipt, err := s.api.GetTransaction(*hash)
	s.NoError(err)
	s.Equal(txstatus.Pending, receipt.Status)

	s.Equal(txtype.MassMigration, receipt.TxType)
}

func TestGetTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionTestSuite))
}
