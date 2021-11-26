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
	s.txsCtx = NewTestTxsContext(s.executionCtx, batchType)
}

func (s *testSuiteWithTxsContext) SetupTestWithConfig(batchType batchtype.BatchType, cfg *config.RollupConfig) {
	s.testSuiteWithExecutionContext.SetupTestWithConfig(cfg)
	s.txsCtx = NewTestTxsContext(s.executionCtx, batchType)
}
