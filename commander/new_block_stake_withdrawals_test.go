package commander

import (
	"context"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	rollupContract "github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type NewBlockLoopSyncStakeWithdrawalsTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd      *Commander
	storage  *st.TestStorage
	client   *eth.TestClient
	cfg      *config.Config
	transfer models.Transfer
	wallets  []bls.Wallet
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DisableSignatures = false

	s.transfer = testutils.MakeTransfer(0, 1, 0, 400)
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.client = newClientWithGenesisStateAndFastBlockFinalization(s.T(), s.storage)

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

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) TestNewBlockLoop_SyncStakeWithdrawals() {
	s.startBlockLoop()
	s.waitForLatestBlockSync()

	s.submitTransferBatchInTransaction(&s.transfer)
	s.waitForLatestBlockSync()

	s.Eventually(s.getStakeWithdrawSendingCondition(0), time.Second, time.Millisecond*50,
		"timeout when waiting for StakeWithdrawEvent")
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) TestNewBlockLoop_SyncStakeWithdrawals_FromScratch() {
	s.startBlockLoop()
	s.waitForLatestBlockSync()

	s.submitTransferBatchInTransaction(&s.transfer)
	s.waitForLatestBlockSync()

	s.Eventually(s.getStakeWithdrawSendingCondition(0), time.Second, time.Millisecond*50,
		"timeout when waiting for StakeWithdrawEvent")

	stopCommander(s.cmd)
	startBlock, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	s.startBlockLoop()
	s.waitForLatestBlockSync()

	s.Never(s.getStakeWithdrawSendingCondition(*startBlock+1), time.Second, time.Millisecond*50,
		"StakeWithdraw must not be sent after restarting the commander")
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) submitTransferBatchInTransaction(tx *models.Transfer) {
	s.runInTransaction(func(txStorage *st.Storage, txsCtx *executor.TxsContext) {
		err := txStorage.AddTransaction(tx)
		s.NoError(err)

		batchData, err := txsCtx.CreateCommitments()
		s.NoError(err)
		s.Len(batchData, 1)

		batch, err := txsCtx.NewPendingBatch(batchtype.Transfer)
		s.NoError(err)
		err = txsCtx.SubmitBatch(batch, batchData)
		s.NoError(err)
		s.client.GetBackend().Commit()
	})
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) runInTransaction(handler func(*st.Storage, *executor.TxsContext)) {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	defer txController.Rollback(nil)

	executionCtx := executor.NewTestExecutionContext(txStorage, s.client.Client, s.cfg.Rollup)
	txsCtx := executor.NewTestTxsContext(executionCtx, batchtype.Transfer)
	handler(txStorage, txsCtx)
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) setAccountsAndChainState() {
	setChainState(s.T(), s.storage)
	setAccountLeaves(s.T(), s.storage.Storage, s.wallets)
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) startBlockLoop() {
	s.cmd.startWorker("", func() error {
		err := s.cmd.newBlockLoop()
		s.NoError(err)
		return nil
	})
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) waitForLatestBlockSync() {
	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	s.Eventually(func() bool {
		syncedBlock, err := s.cmd.storage.GetSyncedBlock()
		s.NoError(err)
		return *syncedBlock >= *latestBlockNumber
	}, time.Hour, 100*time.Millisecond, "timeout when waiting for latest block sync")
}

func (s *NewBlockLoopSyncStakeWithdrawalsTestSuite) getStakeWithdrawSendingCondition(startBlock uint64) func() bool {
	return func() bool {
		it := &rollupContract.StakeWithdrawIterator{}
		latestBlock, err := s.client.GetLatestBlockNumber()
		s.NoError(err)
		s.LessOrEqual(startBlock, *latestBlock)

		err = s.client.FilterLogs(s.client.Rollup.BoundContract, eth.StakeWithdrawEvent, &bind.FilterOpts{
			Start: startBlock,
			End:   latestBlock,
		}, it)
		s.NoError(err)

		for it.Next() {
			if it.Event.Committed == s.client.Blockchain.GetAccount().From {
				return true
			}
		}
		return false
	}
}

func newClientWithGenesisStateAndFastBlockFinalization(t *testing.T, storage *st.TestStorage) *eth.TestClient {
	setStateLeaves(t, storage.Storage)
	genesisRoot, err := storage.StateTree.Root()
	require.NoError(t, err)

	client, err := eth.NewConfiguredTestClient(rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: genesisRoot,
			BlocksToFinalise: models.NewUint256(1),
		},
	}, eth.ClientConfig{})
	require.NoError(t, err)

	return client
}

func TestNewBlockLoopSyncStakeWithdrawalsTestSuite(t *testing.T) {
	suite.Run(t, new(NewBlockLoopSyncStakeWithdrawalsTestSuite))
}
