package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AcceptTransactionsTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.TestStorage
	wallet  *bls.Wallet
	domain  *bls.Domain
}

func (s *AcceptTransactionsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *AcceptTransactionsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.api = &API{
		cfg:               &config.APIConfig{},
		storage:           s.storage.Storage,
		client:            eth.DomainOnlyTestClient,
		commanderMetrics:  metrics.NewCommanderMetrics(),
		disableSignatures: false,
	}

	s.domain, err = s.api.client.GetDomain()
	s.NoError(err)

	s.wallet, err = bls.NewRandomWallet(*s.domain)
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID: 123,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}

	_, err = s.storage.StateTree.Set(1, userState)
	s.NoError(err)
}

func (s *AcceptTransactionsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *AcceptTransactionsTestSuite) TestAcceptTransactions_DisablesSendTransaction() {
	s.api.AcceptTransactions(false)
	s.True(s.api.disableSendTransaction)

	tx := dto.MakeTransaction(transferWithoutSignature)
	transaction, err := s.api.SendTransaction(tx)
	s.Nil(transaction)
	s.ErrorIs(err, ErrSendTransactionDisabled)
}

func (s *AcceptTransactionsTestSuite) TestAcceptTransactions_EnablesSendTransaction() {
	s.api.disableSignatures = true

	s.api.AcceptTransactions(true)
	s.False(s.api.disableSendTransaction)

	tx, err := SignTransfer(s.wallet, transferWithoutSignature)
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(*tx))
	s.NoError(err)
}

func TestAcceptTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptTransactionsTestSuite))
}
