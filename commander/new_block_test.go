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
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type NewBlockLoopTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd      *Commander
	storage  *st.TestStorage
	client   *eth.TestClient
	cfg      *config.Config
	transfer models.Transfer
	wallets  []bls.Wallet
}

func (s *NewBlockLoopTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DisableSignatures = false

	s.transfer = testutils.MakeTransfer(0, 1, 0, 400)
}

func (s *NewBlockLoopTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.client = newClientWithGenesisState(s.T(), s.storage)

	s.cmd = NewCommander(s.cfg, s.client.Blockchain)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.storage.Storage
	s.cmd.metrics = metrics.NewCommanderMetrics()
	s.cmd.workersContext, s.cmd.stopWorkersContext = context.WithCancel(context.Background())

	err = s.cmd.addGenesisBatch()
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)
	s.setAccountsAndChainState()
	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
}

func (s *NewBlockLoopTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_StartsRollupLoop() {
	s.startBlockLoop()

	s.Eventually(func() bool {
		return s.cmd.rollupLoopRunning
	}, 1*time.Second, 100*time.Millisecond)

	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	blockNumber := s.cmd.storage.GetLatestBlockNumber()
	s.Equal(*latestBlockNumber, uint64(blockNumber))
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_SyncsAccountsAndBatchesAndTokensAddedBeforeStartup() {
	accounts := []models.AccountLeaf{
		{PublicKey: *s.wallets[0].PublicKey()},
		{PublicKey: *s.wallets[1].PublicKey()},
	}
	s.registerAccounts(accounts)
	s.submitTransferBatchInTransaction(&s.transfer)
	tokenID := *RegisterSingleToken(s.Assertions, s.client, s.client.ExampleTokenAddress)

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
	s.Len(batches, 2)

	syncedToken, err := s.cmd.storage.GetRegisteredToken(models.MakeUint256(0))
	s.NoError(err)
	s.Equal(s.client.ExampleTokenAddress, syncedToken.Contract)
	s.Equal(tokenID, syncedToken.ID)
}

func (s *NewBlockLoopTestSuite) TestNewBlockLoop_SyncsAccountsAndBatchesAndTokensAddedWhileRunning() {
	s.startBlockLoop()
	s.waitForLatestBlockSync()

	accounts := []models.AccountLeaf{
		{PublicKey: *s.wallets[0].PublicKey()},
		{PublicKey: *s.wallets[1].PublicKey()},
	}
	s.registerAccounts(accounts)
	s.submitTransferBatchInTransaction(&s.transfer)
	tokenID := *RegisterSingleToken(s.Assertions, s.client, s.client.ExampleTokenAddress)

	s.waitForLatestBlockSync()

	for i := range accounts {
		userAccounts, err := s.cmd.storage.AccountTree.Leaves(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Equal(accounts[i], userAccounts[0])
	}

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	syncedToken, err := s.cmd.storage.GetRegisteredToken(models.MakeUint256(0))
	s.NoError(err)
	s.Equal(s.client.ExampleTokenAddress, syncedToken.Contract)
	s.Equal(tokenID, syncedToken.ID)
}

func (s *NewBlockLoopTestSuite) startBlockLoop() {
	s.cmd.startWorker("", func() error {
		err := s.cmd.newBlockLoop()
		s.NoError(err)
		return nil
	})
}

func (s *NewBlockLoopTestSuite) registerAccounts(accounts []models.AccountLeaf) {
	for i := range accounts {
		pubKeyID, err := s.client.RegisterAccountAndWait(&accounts[i].PublicKey)
		s.NoError(err)
		accounts[i].PubKeyID = *pubKeyID
	}
}

func (s *NewBlockLoopTestSuite) submitTransferBatchInTransaction(tx *models.Transfer) {
	s.runInTransaction(func(txStorage *st.Storage, txsCtx *executor.TxsContext) {
		err := txStorage.AddTransfer(tx)
		s.NoError(err)

		batchData, err := txsCtx.CreateCommitments()
		s.NoError(err)
		s.Len(batchData.Commitments(), 1)

		batch, err := txsCtx.NewPendingBatch(batchtype.Transfer)
		s.NoError(err)
		err = txsCtx.SubmitBatch(batch, batchData)
		s.NoError(err)
		s.client.GetBackend().Commit()
	})
}

func (s *NewBlockLoopTestSuite) runInTransaction(handler func(*st.Storage, *executor.TxsContext)) {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	defer txController.Rollback(nil)

	executionCtx := executor.NewTestExecutionContext(txStorage, s.client.Client, s.cfg.Rollup)
	txsCtx := executor.NewTestTxsContext(executionCtx, batchtype.Transfer)
	handler(txStorage, txsCtx)
}

func (s *NewBlockLoopTestSuite) waitForLatestBlockSync() {
	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	s.Eventually(func() bool {
		syncedBlock, err := s.cmd.storage.GetSyncedBlock()
		s.NoError(err)
		return *syncedBlock >= *latestBlockNumber
	}, time.Second, 100*time.Millisecond, "timeout when waiting for latest block sync")
}

func (s *NewBlockLoopTestSuite) setAccountsAndChainState() {
	err := s.storage.SetChainState(&models.ChainState{
		ChainID:     models.MakeUint256(1337),
		SyncedBlock: 0,
	})
	s.NoError(err)

	setAccountLeaves(s.T(), s.storage.Storage, s.wallets)
}

func signTransfer(t *testing.T, wallet *bls.Wallet, transfer *models.Transfer) {
	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	require.NoError(t, err)
	signature, err := wallet.Sign(encodedTransfer)
	require.NoError(t, err)
	transfer.Signature = *signature.ModelsSignature()
}

func setAccountLeaves(t *testing.T, storage *st.Storage, wallets []bls.Wallet) {
	err := storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: *wallets[0].PublicKey(),
	})
	require.NoError(t, err)

	err = storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: *wallets[1].PublicKey(),
	})
	require.NoError(t, err)
}

func setStateLeaves(t *testing.T, storage *st.Storage) {
	_, err := storage.StateTree.Set(0, &models.UserState{
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

	_, err = storage.StateTree.Set(2, &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(2),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	require.NoError(t, err)
}

func stopCommander(cmd *Commander) {
	if !cmd.isRunning {
		return
	}
	cmd.stopWorkersContext()
	cmd.workersWaitGroup.Wait()
}

func TestNewBlockLoopTestSuite(t *testing.T) {
	suite.Run(t, new(NewBlockLoopTestSuite))
}
