package commander

import (
	"context"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
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
	teardown   func() error
}

func (s *NewBlockLoopTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       1,
	}
}

func (s *NewBlockLoopTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)

	s.cmd = NewCommander(config.GetTestConfig())
	s.cmd.client = s.testClient.Client
	s.cmd.storage = testStorage.Storage
	s.cmd.stopChannel = make(chan bool)

	seedDB(s.T(), testStorage.Storage, st.NewStateTree(testStorage.Storage))
}

func (s *NewBlockLoopTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_StartsRollupLoop() {
	s.startBlockLoop()
	defer s.stopCommander()

	s.Eventually(func() bool {
		return s.cmd.rollupLoopRunning
	}, 1*time.Second, 100*time.Millisecond)
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_SyncBeforeLoopStarted() {
	accounts := []models.Account{
		{PublicKey: models.PublicKey{1, 2, 3}},
		{PublicKey: models.PublicKey{2, 3, 4}},
	}
	s.registerAccounts(accounts)

	s.startBlockLoop()
	defer s.stopCommander()
	s.testClient.Commit()

	// TODO: change to eventually
	for i := range accounts {
		userAccounts, err := s.cmd.storage.GetAccounts(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Equal(accounts[i], userAccounts[0])
	}
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_SyncDuringLoop() {
	s.startBlockLoop()
	defer s.stopCommander()

	accounts := []models.Account{
		{PublicKey: models.PublicKey{1, 2, 3}},
		{PublicKey: models.PublicKey{2, 3, 4}},
	}
	s.registerAccounts(accounts)
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID: 1,
	}
	s.createAndSubmitTransferBatch(&tx)
	s.testClient.Commit()

	// TODO: change to eventually
	for i := range accounts {
		userAccounts, err := s.cmd.storage.GetAccounts(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Equal(accounts[i], userAccounts[0])
	}
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
	// TODO: check if commander exited
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
	defer func() {
		err = transactionExecutor.Commit()
		s.NoError(err)
	}()
	commitments, err := transactionExecutor.createTransferCommitments([]models.Transfer{*tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = transactionExecutor.submitBatch(txtype.Transfer, commitments)
	s.NoError(err)
}

func TestNewBlockLoopTestSuite(t *testing.T) {
	suite.Run(t, new(NewBlockLoopTestSuite))
}
