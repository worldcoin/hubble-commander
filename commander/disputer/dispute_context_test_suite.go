package disputer

import (
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuiteWithDisputeContext struct {
	*require.Assertions
	suite.Suite
	storage      *st.TestStorage
	txController *db.TxController
	cfg          *config.RollupConfig
	client       *eth.TestClient
	executionCtx *executor.ExecutionContext
	rollupCtx    *executor.RollupContext
	syncCtx      *syncer.SyncContext
	disputeCtx   *DisputeContext
}

func (s *TestSuiteWithDisputeContext) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TestSuiteWithDisputeContext) SetupTest(batchType batchtype.BatchType) {
	s.SetupTestWithConfig(batchType, config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DisableSignatures:      false,
	})
}

func (s *TestSuiteWithDisputeContext) SetupTestWithConfig(batchType batchtype.BatchType, cfg config.RollupConfig) {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.cfg = &cfg

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, batchType)
}

func (s *TestSuiteWithDisputeContext) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TestSuiteWithDisputeContext) newContexts(
	storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, batchType batchtype.BatchType,
) {
	s.executionCtx = executor.NewTestExecutionContext(storage, s.client.Client, s.cfg)
	s.rollupCtx = executor.NewTestRollupContext(s.executionCtx, batchType)
	s.syncCtx = syncer.NewTestSyncContext(storage, client, cfg, batchType)
	s.disputeCtx = NewDisputeContext(storage, s.client.Client)
}

func (s *TestSuiteWithDisputeContext) beginTransaction() {
	txController, txStorage, err := s.storage.BeginTransaction(st.TxOptions{})
	s.NoError(err)
	s.txController = txController
	s.newContexts(txStorage, s.client.Client, s.cfg, s.rollupCtx.BatchType)
}

func (s *TestSuiteWithDisputeContext) commitTransaction() {
	err := s.txController.Commit()
	s.NoError(err)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, s.rollupCtx.BatchType)
}

func (s *TestSuiteWithDisputeContext) rollback() {
	s.txController.Rollback(nil)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, s.rollupCtx.BatchType)
}

func (s *TestSuiteWithDisputeContext) submitTransferBatch(tx *models.Transfer) *models.Batch {
	pendingBatch, commitments := s.createTransferBatch(tx)

	err := s.rollupCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *TestSuiteWithDisputeContext) createTransferBatch(tx *models.Transfer) (*models.Batch, []models.Commitment) {
	err := s.disputeCtx.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := s.rollupCtx.NewPendingBatch(batchtype.Transfer)
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)

	return pendingBatch, commitments
}

func (s *TestSuiteWithDisputeContext) submitC2TBatch(tx *models.Create2Transfer) {
	pendingBatch, commitments := s.createC2TBatch(tx)

	err := s.rollupCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
}

func (s *TestSuiteWithDisputeContext) createC2TBatch(tx *models.Create2Transfer) (*models.Batch, []models.Commitment) {
	err := s.disputeCtx.storage.AddCreate2Transfer(tx)
	s.NoError(err)

	pendingBatch, err := s.rollupCtx.NewPendingBatch(batchtype.Create2Transfer)
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)
	return pendingBatch, commitments
}
