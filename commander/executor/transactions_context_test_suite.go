package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

type testSuiteWithTxsContext struct {
	testSuiteWithExecutionContext
	txsCtx *TxsContext
}

func (s *testSuiteWithTxsContext) SetupTest(batchType batchtype.BatchType) {
	s.testSuiteWithExecutionContext.SetupTest()
	s.newTestTxsContext(batchType)
}

func (s *testSuiteWithTxsContext) SetupTestWithConfig(batchType batchtype.BatchType, cfg *config.RollupConfig) {
	s.testSuiteWithExecutionContext.SetupTestWithConfig(cfg)
	s.newTestTxsContext(batchType)
}

// AcceptNewConfig testify does not support parameterized test fixtures and propagators are not in
// fashion so when a test changes the RollupConfig it must also redo the relevant parts of setup
func (s *testSuiteWithTxsContext) AcceptNewConfig() {
	batchType := s.txsCtx.BatchType
	s.executionCtx = NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.newTestTxsContext(batchType)
}

func (s *testSuiteWithTxsContext) newTestTxsContext(batchType batchtype.BatchType) {
	var err error
	s.txsCtx, err = NewTestTxsContext(s.executionCtx, batchType)
	s.NoError(err)
}
