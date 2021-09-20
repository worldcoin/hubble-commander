package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

type TestSuiteWithSyncContext struct {
	TestSuiteWithExecutionContext
	rollupCtx *RollupContext
	syncCtx   *SyncContext
}

func (s *TestSuiteWithSyncContext) SetupTest(batchType txtype.TransactionType) {
	s.TestSuiteWithExecutionContext.SetupTest()
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
	s.syncCtx = NewTestSyncContext(s.executionCtx, batchType)
}

func (s *TestSuiteWithSyncContext) SetupTestWithConfig(batchType txtype.TransactionType, cfg config.RollupConfig) {
	s.TestSuiteWithExecutionContext.SetupTestWithConfig(cfg)
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
	s.syncCtx = NewTestSyncContext(s.executionCtx, batchType)
}
