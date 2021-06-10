package commander

import (
	"context"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type NewBlockLoopTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd        *Commander
	testClient *eth.TestClient
	cfg        *config.RollupConfig
	transfer   models.Transfer
	teardown   func() error
}

func (s *NewBlockLoopTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       1,
	}

	s.transfer = models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   s.mockSignature(),
		},
		ToStateID: 1,
	}
}

func (s *NewBlockLoopTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.cmd.storage = testStorage.Storage

	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	s.cmd.client = s.testClient.Client

	s.cmd = NewCommander(config.GetTestConfig())
	s.cmd.stopChannel = make(chan bool)

	seedDB(s.T(), testStorage.Storage, st.NewStateTree(testStorage.Storage))
}

func (s *NewBlockLoopTestSuite) TearDownTest() {
	s.stopCommander()
	err := s.teardown()
	s.NoError(err)
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_StartsRollupLoop() {
	s.startBlockLoop()

	s.Eventually(func() bool {
		return s.cmd.rollupLoopRunning
	}, 1*time.Second, 100*time.Millisecond)

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	blockNumber := s.cmd.storage.GetLatestBlockNumber()
	s.Equal(*latestBlockNumber, uint64(blockNumber))
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_SyncsAccountsAndBatchesAddedBeforeStartup() {
	accounts := []models.Account{
		{PublicKey: models.PublicKey{1, 2, 3}},
		{PublicKey: models.PublicKey{2, 3, 4}},
	}
	s.registerAccounts(accounts)
	s.createAndSubmitTransferBatch(&s.transfer)
	s.testClient.Commit()

	s.startBlockLoop()

	for i := range accounts {
		userAccounts, err := s.cmd.storage.GetAccounts(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Equal(accounts[i], userAccounts[0])
	}

	s.Eventually(func() bool {
		batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
		s.NoError(err)
		return len(batches) == 1
	}, 1*time.Second, 100*time.Millisecond)
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_SyncsAccountsAndBatchesAddedWhileRunning() {
	s.startBlockLoop()

	accounts := []models.Account{
		{PublicKey: models.PublicKey{1, 2, 3}},
		{PublicKey: models.PublicKey{2, 3, 4}},
	}
	s.registerAccounts(accounts)
	s.createAndSubmitTransferBatch(&s.transfer)
	s.testClient.Commit()

	for i := range accounts {
		userAccounts, err := s.cmd.storage.GetAccounts(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Equal(accounts[i], userAccounts[0])
	}

	s.Eventually(func() bool {
		batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
		s.NoError(err)
		return len(batches) == 1
	}, 1*time.Second, 100*time.Millisecond)
}

func (s *NewBlockLoopTestSuite) startBlockLoop() {
	s.cmd.startWorker(func() error {
		err := s.cmd.newBlockLoop()
		s.NoError(err)
		return nil
	})
}

func (s *NewBlockLoopTestSuite) stopCommander() {
	if !s.cmd.IsRunning() {
		return
	}
	close(s.cmd.stopChannel)
	s.cmd.workers.Wait()
}

func (s *NewBlockLoopTestSuite) registerAccounts(accounts []models.Account) {
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	registrations, unsubscribe, err := s.testClient.WatchRegistrations(&bind.WatchOpts{Start: latestBlockNumber})
	s.NoError(err)
	defer unsubscribe()

	for i := range accounts {
		pubKeyID, err := s.testClient.RegisterAccount(&accounts[i].PublicKey, registrations)
		s.NoError(err)
		accounts[i].PubKeyID = *pubKeyID
	}
}

func (s *NewBlockLoopTestSuite) createAndSubmitTransferBatch(tx *models.Transfer) {
	err := s.cmd.storage.AddTransfer(tx)
	s.NoError(err)

	transactionExecutor, err := newTransactionExecutorWithCtx(context.Background(), s.cmd.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	commitments, err := transactionExecutor.createTransferCommitments([]models.Transfer{*tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	_, err = transactionExecutor.submitBatch(txtype.Transfer, commitments)
	s.NoError(err)

	transactionExecutor.Rollback(nil)
}

func (s *NewBlockLoopTestSuite) mockSignature() models.Signature {
	wallet, err := bls.NewRandomWallet(*testDomain)
	s.NoError(err)
	signature, err := wallet.Sign(utils.RandomBytes(4))
	s.NoError(err)
	return *signature.ModelsSignature()
}

func TestNewBlockLoopTestSuite(t *testing.T) {
	suite.Run(t, new(NewBlockLoopTestSuite))
}
