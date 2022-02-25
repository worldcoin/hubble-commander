package syncer

import (
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuiteWithSyncAndRollupContext struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	client  *eth.TestClient
	cfg     *config.RollupConfig
	syncCtx *Context
	txsCtx  *executor.TxsContext
}

func (s *testSuiteWithSyncAndRollupContext) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *testSuiteWithSyncAndRollupContext) SetupTest(batchType batchtype.BatchType) {
	s.SetupTestWithConfig(batchType, &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
	})
}

func (s *testSuiteWithSyncAndRollupContext) SetupTestWithConfig(batchType batchtype.BatchType, cfg *config.RollupConfig) {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.cfg = cfg

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	executionCtx := executor.NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.txsCtx, err = executor.NewTestTxsContext(executionCtx, batchType)
	s.NoError(err)
	s.syncCtx = NewTestContext(s.storage.Storage, s.client.Client, s.cfg, batchType)
}

func (s *testSuiteWithSyncAndRollupContext) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}
