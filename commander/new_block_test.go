package commander

import (
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
	wallets    []bls.Wallet
}

func (s *NewBlockLoopTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       1,
		DevMode:                false,
	}

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
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	err = testStorage.SetChainState(&s.testClient.ChainState)
	s.NoError(err)

	s.cmd = NewCommander(config.GetTestConfig())
	s.cmd.client = s.testClient.Client
	s.cmd.storage = testStorage.Storage
	s.cmd.stopChannel = make(chan bool)

	s.wallets = generateWallets(s.T(), s.testClient.ChainState.Rollup, 2)
	seedDB(s.T(), testStorage.Storage, st.NewStateTree(testStorage.Storage), s.wallets)
	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
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
		{PublicKey: *s.wallets[0].PublicKey()},
		{PublicKey: *s.wallets[1].PublicKey()},
	}
	s.registerAccounts(accounts)
	s.createAndSubmitTransferBatch(&s.transfer)
	s.testClient.Commit()

	s.startBlockLoop()
	s.waitForLatestBlockSync()

	for i := range accounts {
		userAccounts, err := s.cmd.storage.GetAccounts(&accounts[i].PublicKey)
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

	accounts := []models.Account{
		{PublicKey: *s.wallets[0].PublicKey()},
		{PublicKey: *s.wallets[1].PublicKey()},
	}
	s.registerAccounts(accounts)
	s.createAndSubmitTransferBatch(&s.transfer)

	s.testClient.Commit()
	s.waitForLatestBlockSync()

	for i := range accounts {
		userAccounts, err := s.cmd.storage.GetAccounts(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Equal(accounts[i], userAccounts[0])
	}

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
}

func (s *NewBlockLoopTestSuite) TestUnsafeSyncBatches_DoesNotSyncExistingBatchTwice() {
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 1,
	}
	signTransfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	s.createAndSubmitTransferBatch(&tx)
	s.testClient.Commit()

	s.syncAllBlocks()

	tx2 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 0,
	}
	signTransfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	s.createAndSubmitTransferBatch(&tx2)
	s.testClient.Commit()

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)

	s.syncAllBlocks()

	batches, err = s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	state0, err := s.cmd.storage.GetStateLeaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(700), state0.Balance)

	state1, err := s.cmd.storage.GetStateLeaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(300), state1.Balance)
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

	transactionExecutor, err := NewTransactionExecutor(s.cmd.storage, s.testClient.Client, s.cfg, TransactionExecutorOpts{})
	s.NoError(err)

	commitments, err := transactionExecutor.createTransferCommitments([]models.Transfer{*tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	batch, err := transactionExecutor.newPendingBatch(txtype.Transfer)
	s.NoError(err)
	err = transactionExecutor.submitBatch(batch, commitments)
	s.NoError(err)

	transactionExecutor.Rollback(nil)
}

func (s *NewBlockLoopTestSuite) waitForLatestBlockSync() {
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	s.Eventually(func() bool {
		syncedBlock, err := s.cmd.storage.GetSyncedBlock(s.testClient.Client.ChainState.ChainID)
		s.NoError(err)
		return *syncedBlock >= *latestBlockNumber
	}, time.Second, 100*time.Millisecond, "timeout when waiting for latest block sync")
}

func (s *NewBlockLoopTestSuite) syncAllBlocks() {
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	err = s.cmd.unsafeSyncBatches(0, *latestBlockNumber)
	s.NoError(err)
}

func TestNewBlockLoopTestSuite(t *testing.T) {
	suite.Run(t, new(NewBlockLoopTestSuite))
}
