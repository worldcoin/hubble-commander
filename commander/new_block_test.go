package commander

import (
	"context"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
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
	cmd         *Commander
	testStorage *st.TestStorage
	testClient  *eth.TestClient
	cfg         *config.Config
	transfer    models.Transfer
	wallets     []bls.Wallet
}

func (s *NewBlockLoopTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DisableSignatures = false

	s.transfer = models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 1,
	}
}

func (s *NewBlockLoopTestSuite) SetupTest() {
	var err error
	s.testStorage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)

	s.cmd = NewCommander(s.cfg)
	s.cmd.client = s.testClient.Client
	s.cmd.storage = s.testStorage.Storage
	s.cmd.workersContext, s.cmd.stopWorkers = context.WithCancel(context.Background())

	domain, err := s.testClient.GetDomain()
	s.NoError(err)
	s.wallets = generateWallets(s.T(), *domain, 2)
	seedDB(s.T(), s.testStorage.Storage, s.wallets)
	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
}

func (s *NewBlockLoopTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.testClient.Close()
	err := s.testStorage.Teardown()
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
	accounts := []models.AccountLeaf{
		{PublicKey: *s.wallets[0].PublicKey()},
		{PublicKey: *s.wallets[1].PublicKey()},
	}
	s.registerAccounts(accounts)
	s.createAndSubmitTransferBatchInTransaction(&s.transfer)

	s.startBlockLoop()
	s.waitForLatestBlockSync()

	for i := range accounts {
		userAccounts, err := s.cmd.storage.AccountTree.Leaves(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Equal(accounts[i], userAccounts[0])
	}

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_SyncsAccountsAndBatchesAddedWhileRunning() {
	s.startBlockLoop()
	s.waitForLatestBlockSync()

	accounts := []models.AccountLeaf{
		{PublicKey: *s.wallets[0].PublicKey()},
		{PublicKey: *s.wallets[1].PublicKey()},
	}
	s.registerAccounts(accounts)
	s.createAndSubmitTransferBatchInTransaction(&s.transfer)

	s.waitForLatestBlockSync()

	for i := range accounts {
		userAccounts, err := s.cmd.storage.AccountTree.Leaves(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Equal(accounts[i], userAccounts[0])
	}

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
}

func (s *NewBlockLoopTestSuite) startBlockLoop() {
	s.cmd.startWorker(func() error {
		err := s.cmd.newBlockLoop()
		s.NoError(err)
		return nil
	})
}

func (s *NewBlockLoopTestSuite) registerAccounts(accounts []models.AccountLeaf) {
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

func (s *NewBlockLoopTestSuite) createAndSubmitTransferBatchInTransaction(tx *models.Transfer) {
	s.runInTransaction(func(txStorage *st.Storage, txExecutor *executor.TransactionExecutor) {
		_, err := txStorage.AddTransfer(tx)
		s.NoError(err)

		domain, err := s.testClient.GetDomain()
		s.NoError(err)
		commitments, err := txExecutor.CreateTransferCommitments(domain)
		s.NoError(err)
		s.Len(commitments, 1)

		batch, err := txExecutor.NewPendingBatch(txtype.Transfer)
		s.NoError(err)
		err = txExecutor.SubmitBatch(batch, commitments)
		s.NoError(err)
		s.testClient.Commit()
	})
}

func (s *NewBlockLoopTestSuite) runInTransaction(handler func(*st.Storage, *executor.TransactionExecutor)) {
	txController, txStorage, err := s.testStorage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	s.NoError(err)
	defer txController.Rollback(nil)

	txExecutor := executor.NewTestTransactionExecutor(txStorage, s.testClient.Client, s.cfg.Rollup, context.Background())
	handler(txStorage, txExecutor)
}

func (s *NewBlockLoopTestSuite) waitForLatestBlockSync() {
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	s.Eventually(func() bool {
		syncedBlock, err := s.cmd.storage.GetSyncedBlock()
		s.NoError(err)
		return *syncedBlock >= *latestBlockNumber
	}, time.Second, 100*time.Millisecond, "timeout when waiting for latest block sync")
}

func signTransfer(t *testing.T, wallet *bls.Wallet, transfer *models.Transfer) {
	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	require.NoError(t, err)
	signature, err := wallet.Sign(encodedTransfer)
	require.NoError(t, err)
	transfer.Signature = *signature.ModelsSignature()
}

func generateWallets(t *testing.T, domain bls.Domain, walletsAmount int) []bls.Wallet {
	wallets := make([]bls.Wallet, 0, walletsAmount)
	for i := 0; i < walletsAmount; i++ {
		wallet, err := bls.NewRandomWallet(domain)
		require.NoError(t, err)
		wallets = append(wallets, *wallet)
	}
	return wallets
}

func seedDB(t *testing.T, storage *st.Storage, wallets []bls.Wallet) {
	err := storage.SetChainState(&models.ChainState{
		ChainID:     models.MakeUint256(1337),
		SyncedBlock: 0,
	})
	require.NoError(t, err)

	err = storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: *wallets[0].PublicKey(),
	})
	require.NoError(t, err)

	err = storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: *wallets[1].PublicKey(),
	})
	require.NoError(t, err)

	_, err = storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	require.NoError(t, err)

	_, err = storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	require.NoError(t, err)
}

func stopCommander(cmd *Commander) {
	if !cmd.IsRunning() {
		return
	}
	cmd.stopWorkers()
	cmd.workers.Wait()
}

func TestNewBlockLoopTestSuite(t *testing.T) {
	suite.Run(t, new(NewBlockLoopTestSuite))
}
