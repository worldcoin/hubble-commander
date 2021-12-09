package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreateCommitmentsTestSuite struct {
	testSuiteWithTxsContext
	wallets []bls.Wallet
}

func (s *CreateCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CreateCommitmentsTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTestWithConfig(batchtype.Transfer, &config.RollupConfig{
		MinTxsPerCommitment:    2,
		MaxTxsPerCommitment:    32,
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
	})

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)

	s.addUserStates()
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNotCreateCommitmentsWithLessTxsThanRequired() {
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorIfCouldNotCreateEnoughCommitments() {
	s.cfg.MinTxsPerCommitment = 1
	s.cfg.MaxTxsPerCommitment = 1
	s.cfg.MinCommitmentsPerBatch = 2

	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_StoresErrorMessagesOfInvalidTransactions() {
	s.cfg.MinTxsPerCommitment = 1

	invalidTransfer := testutils.MakeTransfer(1, 1234, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &invalidTransfer)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)

	s.Len(s.txsCtx.txErrorsToStore, 1)
	s.Equal(invalidTransfer.Hash, s.txsCtx.txErrorsToStore[0].TxHash)
	s.Equal(applier.ErrNonexistentReceiver.Error(), s.txsCtx.txErrorsToStore[0].ErrorMessage)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNotCallRevertToWhenNotNecessary() {
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	preStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)

	postStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.NotEqual(preStateRoot, postStateRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_CallsRevertToWhenNecessary() {
	validTransfers := []models.Transfer{
		testutils.MakeTransfer(1, 2, 0, 100),
		testutils.MakeTransfer(1, 2, 1, 100),
		testutils.MakeTransfer(1, 2, 2, 100),
	}
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)

	// Calculate state root after applying 2 valid transfers
	s.cfg.MinTxsPerCommitment = 2
	s.cfg.MinCommitmentsPerBatch = 1

	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfers[0])
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfers[1])

	tempTxsCtx := NewTxsContext(
		s.txsCtx.storage,
		s.txsCtx.client,
		s.cfg,
		metrics.NewCommanderMetrics(),
		context.Background(),
		batchtype.Transfer,
	)
	batchData, err := tempTxsCtx.CreateCommitments()
	s.NoError(err)
	s.Equal(batchData.Len(), 1)

	expectedPostStateRoot, err := tempTxsCtx.storage.StateTree.Root()
	s.NoError(err)

	tempTxsCtx.Rollback(nil)

	// Do the test
	s.cfg.MinTxsPerCommitment = 2
	s.cfg.MaxTxsPerCommitment = 2
	s.cfg.MinCommitmentsPerBatch = 1

	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfers[2])
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	batchData, err = s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Equal(batchData.Len(), 1)

	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(expectedPostStateRoot, stateRoot)
}

func (s *CreateCommitmentsTestSuite) hashSignAndAddTransfer(wallet *bls.Wallet, transfer *models.Transfer) {
	hash, err := encoder.HashTransfer(transfer)
	s.NoError(err)
	transfer.Hash = *hash

	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	s.NoError(err)
	signature, err := wallet.Sign(encodedTransfer)
	s.NoError(err)
	transfer.Signature = *signature.ModelsSignature()

	err = s.storage.AddTransfer(transfer)
	s.NoError(err)
}

func (s *CreateCommitmentsTestSuite) addUserStates() {
	feeReceiver := &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	_, err := s.storage.StateTree.Set(0, feeReceiver)
	s.NoError(err)

	_, err = s.storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(2, &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestCreateCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateCommitmentsTestSuite))
}
