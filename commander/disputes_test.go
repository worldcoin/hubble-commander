package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputesTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd                 *Commander
	storage             *st.TestStorage
	client              *eth.TestClient
	transactionExecutor *executor.TransactionExecutor
	cfg                 *config.Config
	wallets             []bls.Wallet
}

func (s *DisputesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DevMode = false
}

func (s *DisputesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.cmd = NewCommander(s.cfg)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.storage.Storage
	s.cmd.workersContext, s.cmd.stopWorkers = context.WithCancel(context.Background())

	s.transactionExecutor = executor.NewTestTransactionExecutor(s.storage.Storage, s.client.Client, s.cfg.Rollup, context.Background())

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = generateWallets(s.T(), *domain, 2)
	seedDB(s.T(), s.storage.Storage, s.wallets)
}

func (s *DisputesTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DisputesTestSuite) TestManageRemoteBatchRollback_StopsRollupLoop() {
	rollupCtx, cancel := context.WithCancel(s.cmd.workersContext)
	defer cancel()

	s.cmd.startWorker(func() error {
		return s.cmd.rollupLoop(rollupCtx)
	})

	invalidBatchID := models.NewUint256(2)
	err := s.cmd.manageRemoteBatchRollback(invalidBatchID, cancel)
	s.NoError(err)
	s.False(s.cmd.rollupLoopRunning)
	s.Equal(*invalidBatchID, s.cmd.invalidBatchID)
}

func (s *DisputesTestSuite) TestManageRemoteBatchRollback_RevertsBatches() {
	_, cancel := context.WithCancel(s.cmd.workersContext)
	defer cancel()

	stateRootBefore, err := s.storage.StateTree.Root()
	s.NoError(err)

	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 100),
		testutils.MakeTransfer(0, 1, 1, 50),
	}
	for i := range transfers {
		s.createAndSubmitTransferBatch(&transfers[i])
	}

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	err = s.cmd.manageRemoteBatchRollback(models.NewUint256(1), cancel)
	s.NoError(err)

	batches, err = s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 0)

	stateRootAfter, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(stateRootBefore, stateRootAfter)
}

func (s *DisputesTestSuite) createAndSubmitTransferBatch(tx *models.Transfer) {
	_, err := s.storage.AddTransfer(tx)
	s.NoError(err)

	batch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	commitments, err := s.transactionExecutor.CreateTransferCommitments(domain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = s.transactionExecutor.SubmitBatch(batch, commitments)
	s.NoError(err)
	s.client.Commit()
}

func TestDisputeWatcherTestSuite(t *testing.T) {
	suite.Run(t, new(DisputesTestSuite))
}
