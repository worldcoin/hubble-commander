package disputer

import (
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuiteWithContexts struct {
	*require.Assertions
	suite.Suite
	storage      *st.TestStorage
	txController *db.TxController
	cfg          *config.RollupConfig
	client       *eth.TestClient
	rollupCtx    *executor.RollupContext
	syncCtx      *syncer.Context
	disputeCtx   *Context
}

func (s *testSuiteWithContexts) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *testSuiteWithContexts) SetupTest(batchType batchtype.BatchType) {
	s.SetupTestWithConfig(batchType, config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DisableSignatures:      false,
	})
}

func (s *testSuiteWithContexts) SetupTestWithConfig(batchType batchtype.BatchType, cfg config.RollupConfig) {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.cfg = &cfg

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, batchType)
}

func (s *testSuiteWithContexts) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *testSuiteWithContexts) newContexts(
	storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, batchType batchtype.BatchType,
) {
	executionCtx := executor.NewTestExecutionContext(storage, s.client.Client, s.cfg)
	s.rollupCtx = executor.NewTestRollupContext(executionCtx, batchType)
	s.syncCtx = syncer.NewTestContext(storage, client, cfg, batchType)
	s.disputeCtx = NewContext(storage, s.client.Client)
}

func (s *testSuiteWithContexts) beginTransaction() {
	txController, txStorage, err := s.storage.BeginTransaction(st.TxOptions{})
	s.NoError(err)
	s.txController = txController
	s.newContexts(txStorage, s.client.Client, s.cfg, s.rollupCtx.BatchType)
}

func (s *testSuiteWithContexts) commitTransaction() {
	err := s.txController.Commit()
	s.NoError(err)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, s.rollupCtx.BatchType)
}

func (s *testSuiteWithContexts) rollback() {
	s.txController.Rollback(nil)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, s.rollupCtx.BatchType)
}

func (s *testSuiteWithContexts) submitBatch(tx models.GenericTransaction) *models.Batch {
	pendingBatch, commitments := s.createBatch(tx)

	err := s.rollupCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.GetBackend().Commit()
	return pendingBatch
}

func (s *testSuiteWithContexts) createBatch(tx models.GenericTransaction) (*models.Batch, []models.TxCommitment) {
	if tx.Type() == txtype.Transfer {
		err := s.disputeCtx.storage.AddTransfer(tx.ToTransfer())
		s.NoError(err)
	} else {
		err := s.disputeCtx.storage.AddCreate2Transfer(tx.ToCreate2Transfer())
		s.NoError(err)
	}

	pendingBatch, err := s.rollupCtx.NewPendingBatch(s.rollupCtx.BatchType)
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)

	return pendingBatch, commitments
}
