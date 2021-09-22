package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

type TestSuiteWithSyncContext struct {
	TestSuiteWithExecutionContext
	rollupCtx *RollupContext
	syncCtx   *SyncContext
}

func (s *TestSuiteWithSyncContext) SetupTest(batchType batchtype.BatchType) {
	s.TestSuiteWithExecutionContext.SetupTest()
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
	s.syncCtx = NewTestSyncContext(s.executionCtx, batchType)
}

func (s *TestSuiteWithSyncContext) SetupTestWithConfig(batchType batchtype.BatchType, cfg config.RollupConfig) {
	s.TestSuiteWithExecutionContext.SetupTestWithConfig(cfg)
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
	s.syncCtx = NewTestSyncContext(s.executionCtx, batchType)
}
