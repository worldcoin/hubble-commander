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
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncStakeWithdrawalsTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd      *Commander
	storage  *st.TestStorage
	client   *eth.TestClient
	cfg      *config.Config
	transfer models.Transfer
	wallets  []bls.Wallet
}

func (s *SyncStakeWithdrawalsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DisableSignatures = false

	s.transfer = testutils.MakeTransfer(0, 1, 0, 400)
}

func (s *SyncStakeWithdrawalsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.client = newClientWithGenesisStateAndFastBlockFinalization(s.T(), s.storage)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)
	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)

	s.setupCommander()
}

func (s *SyncStakeWithdrawalsTestSuite) TearDownTest() {
	s.cmd.stopWorkersAndWait()
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SyncStakeWithdrawalsTestSuite) TestNewBlockLoop_WithdrawsStakesAfterBatchesGetFinalised() {
	s.startBlockLoop()
	s.waitForLatestBlockSync()

	s.submitTransferBatchInTransaction(&s.transfer)
	s.waitForLatestBlockSync()

	s.Eventually(
		s.stakeWithdrawalMinedAfterBlock(0),
		time.Second,
		time.Millisecond*50,
		"timeout when waiting for StakeWithdrawEvent",
	)
}

func (s *SyncStakeWithdrawalsTestSuite) TestNewBlockLoop_DoesNotSendStakeWithdrawalsTwiceAfterRunningCommanderFromScratch() {
	s.startBlockLoop()
	s.waitForLatestBlockSync()

	s.submitTransferBatchInTransaction(&s.transfer)
	s.waitForLatestBlockSync()

	s.Eventually(
		s.stakeWithdrawalMinedAfterBlock(0),
		time.Second,
		time.Millisecond*50,
		"timeout when waiting for StakeWithdrawEvent",
	)

	s.cmd.stopWorkersAndWait()

	startBlock, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.setupCommander()
	s.startBlockLoop()
	s.waitForLatestBlockSync()

	s.Never(
		s.stakeWithdrawalMinedAfterBlock(*startBlock+1),
		time.Second,
		time.Millisecond*50,
		"StakeWithdraw must not be sent after restarting the commander",
	)
}

func (s *SyncStakeWithdrawalsTestSuite) submitTransferBatchInTransaction(tx *models.Transfer) {
	s.runInTransaction(func(txStorage *st.Storage, txsCtx *executor.TxsContext) {
		err := txStorage.AddTransaction(tx)
		s.NoError(err)
		err = txStorage.AddMempoolTx(tx)
		s.NoError(err)

		batchData, err := txsCtx.CreateCommitments(context.Background())
		s.NoError(err)
		s.Len(batchData, 1)

		batch, err := txsCtx.NewPendingBatch(batchtype.Transfer)
		s.NoError(err)
		err = txsCtx.SubmitBatch(context.Background(), batch, batchData)
		s.NoError(err)
		s.client.GetBackend().Commit()
	})
}

func (s *SyncStakeWithdrawalsTestSuite) runInTransaction(handler func(*st.Storage, *executor.TxsContext)) {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	defer txController.Rollback(nil)

	executionCtx := executor.NewTestExecutionContext(txStorage, s.client.Client, s.cfg.Rollup)
	txsCtx, err := executor.NewTestTxsContext(executionCtx, batchtype.Transfer)
	s.NoError(err)
	handler(txStorage, txsCtx)
}

func (s *SyncStakeWithdrawalsTestSuite) setAccountsAndChainState() {
	setChainState(s.T(), s.storage)
	setAccountLeaves(s.T(), s.storage.Storage, s.wallets)
}

func (s *SyncStakeWithdrawalsTestSuite) startBlockLoop() {
	s.cmd.startWorker("Test New Block Loop", func() error {
		err := s.cmd.newBlockLoop()
		s.NoError(err)
		return nil
	})
}

func (s *SyncStakeWithdrawalsTestSuite) waitForLatestBlockSync() {
	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	s.Eventually(func() bool {
		syncedBlock, err := s.cmd.storage.GetSyncedBlock()
		s.NoError(err)
		return *syncedBlock >= *latestBlockNumber
	}, time.Hour, 100*time.Millisecond, "timeout when waiting for latest block sync")
}

func (s *SyncStakeWithdrawalsTestSuite) stakeWithdrawalMinedAfterBlock(startBlock uint64) func() bool {
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

		return it.Next()
	}
}

func newClientWithGenesisStateAndFastBlockFinalization(t *testing.T, storage *st.TestStorage) *eth.TestClient {
	genesisRoot, err := storage.StateTree.Root()
	require.NoError(t, err)

	client, err := eth.NewConfiguredTestClient(&rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: genesisRoot,
			BlocksToFinalise: models.NewUint256(1),
		},
	}, &eth.TestClientConfig{})
	require.NoError(t, err)

	return client
}

func (s *SyncStakeWithdrawalsTestSuite) setupCommander() {
	setStateLeaves(s.T(), s.storage.Storage)
	s.cmd = NewCommander(s.cfg, s.client.Blockchain)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.storage.Storage

	err := s.cmd.addGenesisBatch()
	s.NoError(err)

	s.setAccountsAndChainState()
}

func TestNewBlockLoopSyncStakeWithdrawalsTestSuite(t *testing.T) {
	suite.Run(t, new(SyncStakeWithdrawalsTestSuite))
}
