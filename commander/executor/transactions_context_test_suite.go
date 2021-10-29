package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

type testSuiteWithTransactionsContext struct {
	testSuiteWithExecutionContext
	transactionsCtx *TransactionsContext
}

func (s *testSuiteWithTransactionsContext) SetupTest(batchType batchtype.BatchType) {
	s.testSuiteWithExecutionContext.SetupTest()
	s.transactionsCtx = NewTestTransactionsContext(s.executionCtx, batchType)
}

func (s *testSuiteWithTransactionsContext) SetupTestWithConfig(batchType batchtype.BatchType, cfg config.RollupConfig) {
	s.testSuiteWithExecutionContext.SetupTestWithConfig(cfg)
	s.transactionsCtx = NewTestTransactionsContext(s.executionCtx, batchType)
}
