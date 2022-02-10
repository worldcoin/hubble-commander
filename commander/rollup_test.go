package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/commander/tracker"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RollupTestSuite struct {
	*require.Assertions
	suite.Suite
	tracker.TestSuiteWithTxsSending
	testStorage *storage.TestStorage
	testClient  *eth.TestClient
	commander   *Commander
	wallets     []bls.Wallet
}

func (s *RollupTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RollupTestSuite) SetupTest() {
	var err error
	s.testStorage, err = storage.NewTestStorage()
	s.NoError(err)

	s.testClient, err = eth.NewTestClient()
	s.NoError(err)

	s.commander = &Commander{
		cfg: &config.Config{
			Rollup: &config.RollupConfig{
				MinTxsPerCommitment:    1,
				MaxTxsPerCommitment:    1,
				MinCommitmentsPerBatch: 2,
				MaxCommitmentsPerBatch: 32,
			},
		},
		storage:             s.testStorage.Storage,
		client:              s.testClient.Client,
		metrics:             metrics.NewCommanderMetrics(),
		txsTrackingChannels: s.testClient.TxsChannels,
	}

	domain, err := s.testClient.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)

	s.addUserStates()

	s.StartTxsSending(s.commander.txsTrackingChannels.Requests)
}

func (s *RollupTestSuite) TearDownTest() {
	s.StopTxsSending()
	err := s.testStorage.Teardown()
	s.NoError(err)
	s.testClient.Close()
}

func (s *RollupTestSuite) TestRollupLoopIteration_RollbacksStateOnRollupErrorButStoresInvalidTransactionErrorMessages() {
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.setTxHashAndSign(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.setTxHashAndSign(&s.wallets[1], &invalidTransfer)

	err := s.testStorage.AddTransaction(&validTransfer)
	s.NoError(err)
	err = s.testStorage.AddTransaction(&invalidTransfer)
	s.NoError(err)

	preStateRoot, err := s.testStorage.StateTree.Root()
	s.NoError(err)

	currentBatchType := batchtype.Transfer
	err = s.commander.rollupLoopIteration(context.Background(), &currentBatchType)
	s.NoError(err)

	postStateRoot, err := s.testStorage.StateTree.Root()
	s.NoError(err)

	s.Equal(preStateRoot, postStateRoot)

	storedInvalidTransfer, err := s.testStorage.GetTransfer(invalidTransfer.Hash)
	s.NoError(err)
	s.NotNil(storedInvalidTransfer.ErrorMessage)
	s.Equal(applier.ErrBalanceTooLow.Error(), *storedInvalidTransfer.ErrorMessage)
}

func (s *RollupTestSuite) TestRollupLoopIteration_RerunIterationWhenNotEnoughDeposits() {
	s.commander.cfg.Rollup.MinCommitmentsPerBatch = 1
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.setTxHashAndSign(&s.wallets[0], &validTransfer)

	err := s.testStorage.AddTransaction(&validTransfer)
	s.NoError(err)

	currentBatchType := batchtype.Deposit
	err = s.commander.rollupLoopIteration(context.Background(), &currentBatchType)
	s.NoError(err)

	batches, err := s.commander.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(batchtype.Create2Transfer, currentBatchType)
}

func (s *RollupTestSuite) TestRollupLoopIteration_SavesTxErrors() {
	s.commander.cfg.Rollup.MinCommitmentsPerBatch = 1
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.setTxHashAndSign(&s.wallets[0], &validTransfer)

	invalidTransfer := testutils.MakeTransfer(0, 2, 0, 100)
	s.setTxHashAndSign(&s.wallets[0], &validTransfer)

	err := s.testStorage.BatchAddTransfer([]models.Transfer{validTransfer, invalidTransfer})
	s.NoError(err)

	currentBatchType := batchtype.Transfer
	err = s.commander.rollupLoopIteration(context.Background(), &currentBatchType)
	s.NoError(err)

	transfer, err := s.commander.storage.GetTransfer(invalidTransfer.Hash)
	s.NoError(err)
	s.NotNil(transfer.ErrorMessage)
	s.Equal(applier.ErrBalanceTooLow.Error(), *transfer.ErrorMessage)
}

func (s *RollupTestSuite) setTxHashAndSign(wallet *bls.Wallet, transfer *models.Transfer) {
	hash, err := encoder.HashTransfer(transfer)
	s.NoError(err)
	transfer.Hash = *hash

	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	s.NoError(err)
	signature, err := wallet.Sign(encodedTransfer)
	s.NoError(err)
	transfer.Signature = *signature.ModelsSignature()
}

func (s *RollupTestSuite) addUserStates() {
	feeReceiver := &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	_, err := s.testStorage.Storage.StateTree.Set(0, feeReceiver)
	s.NoError(err)

	_, err = s.testStorage.Storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = s.testStorage.Storage.StateTree.Set(2, &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestRollupTestSuite(t *testing.T) {
	suite.Run(t, new(RollupTestSuite))
}
